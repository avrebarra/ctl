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
By adding these:
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/avrebarra/ctl"
)

func main() {
	cpx := ctl.GetGlobal()
	
	// setting config values
	cpx.Set("flags.hello_enabled", true) // no error handling

	// setup handlers
	http.DefaultServeMux.Handle("/ctl/", ctl.MakeHandler(ctl.ConfigHandler{PathPrefix: "/ctl/", Ctl: ctl.GetGlobal()}))
	http.DefaultServeMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// get config value
		flaghello, _ := cpx.Get("flags.hello_enabled").Bool()
		
		// use it
		if flaghello {
			fmt.Fprintf(w, "hello\n")
			return
		}
		fmt.Fprintf(w, "sorry, no hello\n")
	})

	fmt.Println("listening http://localhost:3333...")
	http.ListenAndServe(":3333", http.DefaultServeMux)
}
```

You will have a dynamic configurables plus HTTP control panel:

```sh
$ curl --location --request GET 'localhost:3333/ctl/config'
{"flags.hello_enabled":"true","settings.logging_defaults":"{\"Version\":\"1.0\",\"ClusterID\":\"88888\"}","settings.logging_min_amount":"100000","settings.logging_prefix":"trx_log"}

$ curl --location --request PATCH 'localhost:3333/ctl/config/flags.hello_enabled' \
--header 'Content-Type: application/json' \
--data-raw '{
    "value": "false"
}'
{"key":"flags.transaction_logging_enabled","value":"false"}

```

Available endpoints:
- GET http://localhost:3333/{prefix}/config
- GET http://localhost:3333/{prefix}/config/flags.hello_enabled
- PATCH http://localhost:3333/{prefix}/config/flags.hello_enabled with payload `{ "value":"value to persist AS A STRING" }`

## Other Examples
### Managing values via code
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
}
```

### Referencing and Refreshing Values
```go
package main

import (
	"fmt"

	"github.com/avrebarra/ctl"
)

func main() {
	ref := ctl.GetGlobal().Get("flags.isok") // initial value: none
	fmt.Println(ref.Refresh().Bool())

	ctl.GetGlobal().Set("flags.isok", "true") // change value
	fmt.Println(ref.Refresh().Bool())

	// Output:
	// false value not found
	// true <nil>
}

```

### Replacing the Global Singleton
By default the global instance will be generated with memstore, but you can override it using custom store and options.

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
- [x] Reference values
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
