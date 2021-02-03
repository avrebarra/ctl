package ctl

import (
	"encoding/json"
	"strconv"
)

type Value struct {
	ctlref *Ctl
	key    string
	valsrc string
	err    error
}

func (v Value) Refresh() Value {
	return v.ctlref.Get(v.key)
}

func (v Value) Err() (err error) {
	return v.err
}

// ***

func (v Value) Bool() (val bool, err error) {
	if v.err != nil {
		err = v.err
		return
	}

	return strconv.ParseBool(v.valsrc)
}

func (v Value) Float() (val float64, err error) {
	if v.err != nil {
		err = v.err
		return
	}

	return strconv.ParseFloat(v.valsrc, 64)
}

func (v Value) Int() (val int, err error) {
	err = v.err
	if err != nil {
		return
	}

	return strconv.Atoi(v.valsrc)
}

func (v Value) String() (val string, err error) {
	err = v.err
	if err != nil {
		return
	}

	return v.valsrc, nil
}

func (v Value) Bind(target interface{}) (err error) {
	err = v.err
	if err != nil {
		return
	}

	return json.Unmarshal([]byte(v.valsrc), target)
}
