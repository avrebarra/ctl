package ctl

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Config struct {
	InitialValues map[string]string
}

type Ctl struct {
	config         Config
	valsrc         []byte
	valmap         map[string]string
	subscribersmap map[string][]func(v Value)
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
	}

	ent.Reset()

	return ent, nil
}

func (e *Ctl) Get(key string) (val Value) {
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
		e.valmap[key] = string(btvalue)
		break
	default:
		e.valmap[key] = fmt.Sprint(value)
	}

	// rebuild map
	err := e.rebuildmap()
	if err != nil {
		val.err = err
		return
	}

	// trigger subscribers
	e.triggerSubscribers(key, e.Get(key))

	return
}

func (e *Ctl) Subscribe(key string, fun func(v Value)) (subid string) {
	e.subscribersmap[key] = append(e.subscribersmap[key], fun)
	return
}

func (e *Ctl) Reset() {
	if e.config.InitialValues != nil {
		e.valmap = e.config.InitialValues
		_ = e.rebuildmap() // intendedly thrown

		for k := range e.valmap {
			e.triggerSubscribers(k, e.Get(k))
		}
	}
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
