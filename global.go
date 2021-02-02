package ctl

var globalinst *Ctl

func RegisterGlobal(c *Ctl) {
	globalinst = c
}

func GetGlobal() (c *Ctl) {
	return globalinst
}
