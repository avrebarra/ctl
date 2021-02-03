package ctl

type StoreMem struct {
	strstore string
}

func NewStoreMem() Store {
	return &StoreMem{}
}

func (e *StoreMem) Get() (s string, err error) {
	s = e.strstore
	return
}

func (e *StoreMem) Set(s string) (err error) {
	e.strstore = s
	return
}
