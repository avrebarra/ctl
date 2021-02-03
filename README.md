<div class="info" align="left">
  <h1 class="name">🎛️ ctl</h1>
  Add a control board to golang app/server.
  <br>
  <br>

[![Godoc][godoc-image]][godoc-url]
[![Report][report-image]][report-url]
[![Tests][tests-image]][tests-url]
[![Coverage][coverage-image]][coverage-url]

</div>

## Usage
*PS: Check example_test.go for more up to date examples*

### Basic Usage
```go
package main

import (
	"fmt"

	"github.com/avrebarra/ctl"
)

func main() {
	// setup instance
	xctl, _ := ctl.New(ctl.Config{})

	// read configuration
	fmt.Println(xctl.Get("flags.enable_banner").Bool())
	fmt.Println(xctl.Get("setting.volume").Int())
	fmt.Println(xctl.Get("setting.unknown").String())

	// add configuration
	xctl.Set("flags.enable_debug", true)

	// register as global (optional)
	ctl.RegisterGlobal(xctl)
	fmt.Println(ctl.GetGlobal().Get("flags.enable_banner").Bool())
}
```

### Complex Value (Struct / Map)
```go
package main

import (
	"fmt"

	"github.com/avrebarra/ctl"
)

func main() {
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
```

### Persistence
```go
package main

import (
	"fmt"

	"github.com/avrebarra/ctl"
)

func main() {
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
```

### CTL Endpoint
```go
package main

import (
	"net/http"

	"github.com/avrebarra/ctl"
)

func main() {
	store, _ := ctl.NewStoreFile(ctl.ConfigStoreFile{FilePath: "./fixture/store.json"})
	xctl, _ := ctl.New(ctl.Config{Store: store})

	http.Handle("/ctl", ctl.MakeHandler(ConfigHandler{
		PathPrefix: "/ctl",
		Ctl:        xctl,
	}))

	http.ListenAndServe(":3333", http.DefaultServeMux)
}
```

Available endpoints:
- `GET {prefix}/configs/`
- `GET {prefix}/configs/flags.enable_debug`
- `PUT {prefix}/configs/flags.enable_debug { value:"value" }`


## Milestones
- [x] Value.Float()
- [x] Ctl.Reset()
- [x] Ctl.Subscribe()
- [ ] Ctl.StopSubscribe()
- [x] Persistence
- [ ] Value Encryption
- [x] REST API Handler helper for management

[godoc-image]: https://godoc.org/github.com/avrebarra/ctl?status.svg
[godoc-url]: https://godoc.org/github.com/avrebarra/ctl
[report-image]: https://goreportcard.com/badge/github.com/avrebarra/ctl
[report-url]: https://goreportcard.com/report/github.com/avrebarra/ctl
[tests-image]: https://cloud.drone.io/api/badges/avrebarra/ctl/status.svg
[tests-url]: https://cloud.drone.io/avrebarra/ctl
[coverage-image]: https://codecov.io/gh/avrebarra/ctl/graph/badge.svg
[coverage-url]: https://codecov.io/gh/avrebarra/ctl
[sponsor-image]: https://img.shields.io/badge/github-donate-green.svg