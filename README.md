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
```go
package main

import (
	"fmt"

	"github.com/avrebarra/ctl"
)

func main() {
	// setup instance
	xctl, err := ctl.New(ctl.Config{
		InitialValues: map[string]string{
			"flags.enable_banner": "true",
			"setting.volume":      "5",
		},
	})
	if err != nil {
		panic(err)
	}

	// read configuration
	fmt.Println(xctl.Get("flags.enable_banner").Bool())
	fmt.Println(xctl.Get("setting.volume").Int())
	fmt.Println(xctl.Get("setting.unknown").String())

	// add configuration
	err = xctl.Set("flags.enable_debug", true).Err()
	if err != nil {
		panic(err)
	}

	// register as global (optional)
	ctl.RegisterGlobal(xctl)
	fmt.Println(ctl.GetGlobal().Get("flags.enable_banner").Bool())

	// Output:
	// true <nil>
	// 5 <nil>
	//  value not found
	// false value not found
}
```

## Milestones
- [x] Value.Float()
- [x] Ctl.Reset()
- [x] Ctl.Subscribe()
- [ ] Ctl.StopSubscribe()
- [ ] Persistence
- [ ] Value Encryption
- [ ] REST API Handler helper for management

```
GET /
POST /reset
GET /items/
GET /items/flags.enable_debug
PUT /items/flags.enable_debug
```

[godoc-image]: https://godoc.org/github.com/avrebarra/minimok?status.svg
[godoc-url]: https://godoc.org/github.com/avrebarra/minimok
[report-image]: https://goreportcard.com/badge/github.com/avrebarra/minimok
[report-url]: https://goreportcard.com/report/github.com/avrebarra/minimok
[tests-image]: https://cloud.drone.io/api/badges/avrebarra/minimok/status.svg
[tests-url]: https://cloud.drone.io/avrebarra/minimok
[coverage-image]: https://codecov.io/gh/avrebarra/minimok/graph/badge.svg
[coverage-url]: https://codecov.io/gh/avrebarra/minimok
[sponsor-image]: https://img.shields.io/badge/github-donate-green.svg