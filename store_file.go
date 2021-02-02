package ctl

import (
	"io/ioutil"

	"github.com/go-playground/validator"
)

type ConfigStoreFile struct {
	FilePath string `validate:"required"`
}

type StoreFile struct {
	config ConfigStoreFile
}

func NewStoreFile(cfg ConfigStoreFile) (Store, error) {
	if err := validator.New().Struct(cfg); err != nil {
		return nil, err
	}
	return &StoreFile{config: cfg}, nil
}

func (e *StoreFile) Get() (s string, err error) {
	b, err := ioutil.ReadFile(e.config.FilePath)
	if err != nil {
		return
	}

	s = string(b)

	return
}

func (e *StoreFile) Set(s string) (err error) {
	return ioutil.WriteFile(e.config.FilePath, []byte(s), 0644)
}
