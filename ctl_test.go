package ctl_test

import (
	"testing"

	"github.com/avrebarra/ctl"
	"github.com/stretchr/testify/assert"
)

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
