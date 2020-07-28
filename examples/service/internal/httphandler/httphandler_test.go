package httphandler

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/nexmoinc/gosrvlib/pkg/httpserver/route"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	hh := New(nil)
	require.NotNil(t, hh)
}

func TestHTTPHandler_BindHTTP(t *testing.T) {
	type fields struct {
		service Service
	}
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []route.Route
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTPHandler{
				service: tt.fields.service,
			}
			if got := h.BindHTTP(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BindHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPHandler_handleGenUID(t *testing.T) {
	type fields struct {
		service Service
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTPHandler{
				service: tt.fields.service,
			}
		})
	}
}
