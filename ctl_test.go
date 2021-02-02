package ctl_test

import (
	"fmt"
	"testing"

	"github.com/avrebarra/ctl"
	"github.com/stretchr/testify/assert"
)

func ExampleUsage() {
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
	// true <nil>
}

func TestNew(t *testing.T) {
	type args struct {
		cfg ctl.Config
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				cfg: ctl.Config{},
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "ok/with_initial_values",
			args: args{
				cfg: ctl.Config{
					InitialValues: map[string]string{
						"i1": "data",
						"i2": `{"foal":false}`,
					},
				},
			},
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ctl.New(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err)
				return
			}
			if tt.wantNil {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestCtl_Get(t *testing.T) {
	e, err := ctl.New(ctl.Config{InitialValues: map[string]string{
		"i1": "ok",
	}})
	assert.Nil(t, err)

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantVal bool
		wantErr bool
	}{
		{
			name:    "ok",
			args:    args{key: "i1"},
			wantVal: true,
			wantErr: false,
		},
		{
			name:    "err/not_found",
			args:    args{key: "i2"},
			wantVal: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal := e.Get(tt.args.key)
			if tt.wantVal {
				assert.NotNil(t, gotVal)
			} else {
				assert.Nil(t, gotVal)
			}

			if tt.wantErr {
				assert.NotNil(t, gotVal.Err())
			} else {
				assert.Nil(t, gotVal.Err())
			}
		})
	}
}
