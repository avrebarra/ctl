package ctl

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type Config struct {
	Store       Store
	RefreshRate time.Duration
}

type Ctl struct {
	config         Config
	valsrc         []byte
	valmap         map[string]string
	subscribersmap map[string][]func(v Value)

	nextrefresh time.Time
}

func New(cfg Config) (*Ctl, error) {
	// if err := validator.New().Struct(cfg); err != nil {
	// 	return nil, err
	// }

	ent := &Ctl{
		config:         cfg,
		valsrc:         []byte{},
		valmap:         map[string]string{},
		subscribersmap: map[string][]func(v Value){},
		nextrefresh:    time.Now().Add(cfg.RefreshRate),
	}

	if err := ent.storerefresh(); err != nil {
		return nil, err
	}

	return ent, nil
}

func (e *Ctl) List() (lis map[string]string) {
	e.storerefresh()

	return e.valmap
}

func (e *Ctl) Get(key string) (val Value) {
	e.storerefresh()

	strval, ok := e.valmap[key]
	if !ok {
		val.err = fmt.Errorf("value not found")
		return
	}

	val = Value{valsrc: strval}

	return
}

func (e *Ctl) Set(key string, value interface{}) (val Value) {
	// serialize value
	switch reflect.ValueOf(value).Kind() {
	case reflect.Struct, reflect.Array, reflect.Map:
		btvalue, err := json.Marshal(value)
		if err != nil {
			val.err = err
			return
		}
		val = Value{valsrc: string(btvalue)}
		break
	default:
		val = Value{valsrc: fmt.Sprint(value)}
	}

	e.valmap[key] = val.valsrc

	// rebuild map
	err := e.rebuildmap()
	if err != nil {
		val.err = err
		return
	}

	// trigger subscribers
	e.triggerSubscribers(key, e.Get(key))

	// persist
	err = e.storepersist()
	if err != nil {
		val.err = err
		return
	}

	return
}

func (e *Ctl) Subscribe(key string, fun func(v Value)) (subid string) {
	e.subscribersmap[key] = append(e.subscribersmap[key], fun)
	return
}

// ***

func (e *Ctl) storerefresh() (err error) {
	if e.config.Store == nil {
		return
	}
	if time.Now().Before(e.nextrefresh) {
		return
	}

	// get stored value
	var val string
	val, err = e.config.Store.Get()
	if err != nil {
		return
	}

	// refresh valmap
	e.valsrc = []byte(val)
	err = json.Unmarshal(e.valsrc, &e.valmap)
	if err != nil {
		return
	}

	e.nextrefresh = time.Now().Add(e.config.RefreshRate)

	return
}

func (e *Ctl) storepersist() (err error) {
	if e.config.Store == nil {
		return
	}

	return e.config.Store.Set(string(e.valsrc))
}

func (e *Ctl) rebuildmap() (err error) {
	newvalsrc, err := json.Marshal(e.valmap)
	if err != nil {
		return
	}

	e.valsrc = newvalsrc
	return
}

func (e *Ctl) triggerSubscribers(key string, val Value) (err error) {
	subfuns, ok := e.subscribersmap[key]
	if !ok {
		return
	}

	for _, subfun := range subfuns {
		subfun(val)
	}

	return
}
