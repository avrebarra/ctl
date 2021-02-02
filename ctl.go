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
	config Config
	valsrc []byte
	valmap map[string]string
}

func New(cfg Config) (*Ctl, error) {
	// if err := validator.New().Struct(cfg); err != nil {
	// 	return nil, err
	// }

	ent := &Ctl{
		config: cfg,
		valsrc: []byte{},
		valmap: map[string]string{},
	}

	if cfg.InitialValues != nil {
		ent.valmap = cfg.InitialValues
		_ = ent.rebuildmap() // intendedly thrown
	}

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

	return
}

func (e *Ctl) Subscribe(key string, fun func(v Value)) (subid string) {
	return
}

func (e *Ctl) rebuildmap() (err error) {
	newvalsrc, err := json.Marshal(e.valmap)
	if err != nil {
		return
	}

	e.valsrc = newvalsrc
	return
}
