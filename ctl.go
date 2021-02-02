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

	return
}

func (e *Ctl) Subscribe(key string, fun func(v Value)) (subid string) {
	return
}

func (e *Ctl) Reset() {
	if e.config.InitialValues != nil {
		e.valmap = e.config.InitialValues
		_ = e.rebuildmap() // intendedly thrown
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
