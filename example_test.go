package ctl_test

import (
	"fmt"

	"github.com/avrebarra/ctl"
)

func ExampleBasic() {
	// setup instance
	xctl, err := ctl.New(ctl.Config{})
	if err != nil {
		panic(err)
	}

	// change/add new configuration
	err = xctl.Set("flags.enable_debug", true).Err()
	if err != nil {
		panic(err)
	}
	xctl.Set("flags.enable_banner", true)
	xctl.Set("setting.volume", 5)
	xctl.Set("setting.redeem_rate", .58)

	// read configuration
	fmt.Println(xctl.Get("flags.enable_banner").Bool())
	fmt.Println(xctl.Get("setting.volume").Int())
	fmt.Println(xctl.Get("setting.redeem_rate").Float())
	fmt.Println(xctl.Get("setting.unknown").String())

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
	})

	fmt.Println(xctl.Get("flags.previously_persisted").Bool())
	fmt.Println(xctl.Set("flags.another_example_flag", true).Bool())

	// Output:
	// true <nil>
	// true <nil>
}

func ExampleComplexObject() {
	type objstr struct {
		Data1     string
		Data2     bool
		SubStruct struct {
			Data1 string
			Data2 bool
		}
	}

	// setup instance
	xctl, _ := ctl.New(ctl.Config{})

	// add configuration
	xctl.Set("my_settings.complex_object", objstr{
		Data1: "something",
		Data2: true,
		SubStruct: struct {
			Data1 string
			Data2 bool
		}{
			Data1: "awyeah",
		},
	})

	// read configuration
	got := objstr{}
	err := xctl.Get("my_settings.complex_object").Bind(&got)
	if err != nil {
		panic(err)
	}

	fmt.Println(got.Data1)
	fmt.Println(got.Data2)
	fmt.Println(got.SubStruct.Data1)

	// Output:
	// something
	// true
	// awyeah
}
