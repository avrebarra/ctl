package ctl

type Store interface {
	Get() (s string, err error)
	Set(s string) (err error)
}
