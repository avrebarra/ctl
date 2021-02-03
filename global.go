package ctl

var globalinst *Ctl

func RegisterGlobal(c *Ctl) {
	globalinst = c
}

func GetGlobal() (c *Ctl) {
	if globalinst == nil {
		globalinst, _ = New(Config{Store: NewStoreMem(), RefreshRate: 0})
	}

	return globalinst
}
