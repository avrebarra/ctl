package ctl

import (
	"fmt"
	"testing"
)

func TestValue_Float(t *testing.T) {
	type fields struct {
		valsrc string
		err    error
	}
	tests := []struct {
		name    string
		fields  fields
		wantVal float64
		wantErr bool
	}{
		{
			name: "ok/variant_normal",
			fields: fields{
				valsrc: "0.324",
				err:    nil,
			},
			wantVal: 0.324,
			wantErr: false,
		},
		{
			name: "ok/variant_compact",
			fields: fields{
				valsrc: ".324",
				err:    nil,
			},
			wantVal: 0.324,
			wantErr: false,
		},
		{
			name: "ok/variant_long",
			fields: fields{
				valsrc: ".3248128881203",
				err:    nil,
			},
			wantVal: 0.3248128881203,
			wantErr: false,
		},
		{
			name: "error/cascaded",
			fields: fields{
				valsrc: "0.324",
				err:    fmt.Errorf("random error"),
			},
			wantVal: 0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Value{
				valsrc: tt.fields.valsrc,
				err:    tt.fields.err,
			}
			gotVal, err := v.Float()
			if (err != nil) != tt.wantErr {
				t.Errorf("Value.Float() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Value.Float() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}
