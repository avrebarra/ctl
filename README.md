<div class="info" align="left">
  <h1 class="name">🎛️ ctl</h1>
  Quickly add dynamic configurations and control panel to your app & server.
  <br>
  <br>

[![Godoc][godoc-image]][godoc-url]
[![Report][report-image]][report-url]
[![Tests][tests-image]][tests-url]
[![Coverage][coverage-image]][coverage-url]

</div>

## Usage

### Setup and managing values
*Note: It's recommended to specify a centralized storage. By doing so, multiple instances of same service could make use of shared/synchronized dynamic configs. You can also define your own store for db/redis/consul etc by implementing `Store` interface*

```go
package main

import (
	"fmt"

	"github.com/avrebarra/ctl"
)

func main() {
	// setup instance with file storage
	store, _ := ctl.NewStoreFile(ctl.ConfigStoreFile{FilePath: "fixture/store.json"})
	cpx, _ := ctl.New(ctl.Config{
		Store:       store,
		RefreshRate: 10 * time.Second, // how often to refetch data from store
	})

	// setting configurations
	cpx.Set("flags.logging_enabled", true) // no error handling
	cpx.Set("settings.logging_prefix", "trx_log")
	cpx.Set("settings.logging_defaults", DataField{Version: "1.0", ClusterID: "88888"})

	// handling for errors
	err := cpx.Set("settings.logging_min_amount", 100000).Err()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// getting config and assert as multiple types (boolean, float, string, object)
	flagEnableBanner, _ := cpx.Get("flags.enable_banner").Bool()
	confMinAmt, _ := cpx.Get("settings.logging_min_amount").Int()

	// binding value to object
	datafield := DataField{}
	_ = cpx.Get("settings.logging_defaults").Bind(&datafield)


	// register as global for centralized access on runtime
	ctl.RegisterGlobal(cpx)
	flagEnableBanner, _ = ctl.GetGlobal().Get("flags.enable_banner").Bool()

	fmt.Println("values:", flagEnableBanner, confMinAmt, datafield)
}
```

### Setup HTTP based Control Panel
Ctl also support attaching some endpoints to your HTTP server to enable value management via HTTP request:

```go
package main

import (
	"fmt"

	"github.com/avrebarra/ctl"
)

func main() {
	store, _ := ctl.NewStoreFile(ctl.ConfigStoreFile{FilePath: "fixture/store.json"})
	cpx, _ := ctl.New(ctl.Config{
		Store:       store,
		RefreshRate: 10 * time.Second,
	})

	http.DefaultServeMux.Handle("/ctl/", ctl.MakeHandler(ctl.ConfigHandler{
		PathPrefix: "/ctl/",
		Ctl:        cpx,
	}))

	fmt.Println("listening http://localhost:3333...")
	fmt.Println("     to see ctl config listing, visit http://localhost:3333/ctl/config")
	fmt.Println("     get and update individual config using GET/POST http://localhost:3333/ctl/config/{keys}")
	http.ListenAndServe(":3333", http.DefaultServeMux)
}
```

Available endpoints:
- GET http://localhost:3333/{prefix}/config
- GET http://localhost:3333/{prefix}/config/flags.enable_debug
- PUT http://localhost:3333/{prefix}/config/flags.enable_debug with payload `{ "value":"value to persist in string" }`

### Register Global Instance
```go
package main

import (
	"fmt"

	"github.com/avrebarra/ctl"
)

func main() {
	store, _ := ctl.NewStoreFile(ctl.ConfigStoreFile{FilePath: "fixture/store.json"})
	cpx, _ := ctl.New(ctl.Config{
		Store:       store,
		RefreshRate: 10 * time.Second,
	})

	ctl.RegisterGlobal(cpx)

	flagEnableBanner, _ = ctl.GetGlobal().Get("flags.enable_banner").Bool()
	fmt.Println("values:", flagEnableBanner)
}
```

## Milestones
- [x] Value.Float()
- [x] Ctl.Reset()
- [x] Ctl.Subscribe()
- [ ] Ctl.StopSubscribe()
- [x] Persistence
- [ ] Value Encryption
- [ ] Pointer values
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