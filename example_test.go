package ctl_test

import (
	"fmt"

	"github.com/avrebarra/ctl"
)

func ExampleBasic() {
	// setup instance
	xctl, err := ctl.New(ctl.Config{
		InitialValues: map[string]string{
			"flags.enable_banner": "true",
			"setting.volume":      "5",
			"setting.redeem_rate": ".58",
		},
	})
	if err != nil {
		panic(err)
	}

	// read configuration
	fmt.Println(xctl.Get("flags.enable_banner").Bool())
	fmt.Println(xctl.Get("setting.volume").Int())
	fmt.Println(xctl.Get("setting.redeem_rate").Float())
	fmt.Println(xctl.Get("setting.unknown").String())

	// change/add new configuration
	err = xctl.Set("flags.enable_debug", true).Err()
	if err != nil {
		panic(err)
	}

	// subscribe to changes
	xctl.Subscribe("setting.volume", func(v ctl.Value) {
		fmt.Println("changed: setting.volume")
	})
	xctl.Set("setting.volume", 6)

	// register as global (optional)
	ctl.RegisterGlobal(xctl)
	fmt.Println(ctl.GetGlobal().Get("flags.enable_banner").Bool())

	// Output:
	// true <nil>
	// 5 <nil>
	// 0.58 <nil>
	//  value not found
	// changed: setting.volume
	// true <nil>
}

func ExamplePersistence() {
	// setup store
	fstore, _ := ctl.NewStoreFile(ctl.ConfigStoreFile{
		FilePath: "fixture/store.json",
	})

	// setup instance
	xctl, _ := ctl.New(ctl.Config{
		Store: fstore,
		InitialValues: map[string]string{
			"flags.enable_banner": "true",
			"setting.volume":      "5",
			"setting.redeem_rate": ".58",
		},
	})

	// fmt.Println(xctl.Get("flags.previously_persisted").Bool())
	fmt.Println(xctl.Set("flags.another_example_flag", true).Bool())

	// Output:
	// true <nil>
	// true <nil>
}
