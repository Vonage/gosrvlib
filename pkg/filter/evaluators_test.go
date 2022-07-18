package filter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExact_Evaluate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		ref   interface{}
		value interface{}
		want  bool
	}{
		{
			name:  "true - int / int",
			ref:   42,
			value: 42,
			want:  true,
		},
		{
			name:  "true - float64 / int",
			ref:   42.0,
			value: 42,
			want:  true,
		},
		{
			name:  "true - int / float64",
			ref:   42,
			value: 42.0,
			want:  true,
		},
		{
			name:  "true - float64 / float64",
			ref:   42.0,
			value: 42.0,
			want:  true,
		},
		{
			name:  "true - int8 / float64",
			ref:   int8(42),
			value: 42.0,
			want:  true,
		},
		{
			name:  "true - int16 / float64",
			ref:   int16(42),
			value: 42.0,
			want:  true,
		},
		{
			name:  "true - int32 / float64",
			ref:   int32(42),
			value: 42.0,
			want:  true,
		},
		{
			name:  "true - int64 / float64",
			ref:   int64(42),
			value: 42.0,
			want:  true,
		},
		{
			name:  "true - uint / float64",
			ref:   uint(42),
			value: 42.0,
			want:  true,
		},
		{
			name:  "true - uint8 / float64",
			ref:   uint8(42),
			value: 42.0,
			want:  true,
		},
		{
			name:  "true - uint16 / float64",
			ref:   uint16(42),
			value: 42.0,
			want:  true,
		},
		{
			name:  "true - uint32 / float64",
			ref:   uint32(42),
			value: 42.0,
			want:  true,
		},
		{
			name:  "true - uint64 / float64",
			ref:   uint64(42),
			value: 42.0,
			want:  true,
		},
		{
			name:  "true - float32 / float64",
			ref:   float32(42),
			value: 42.0,
			want:  true,
		},
		{
			name:  "false - int / int",
			ref:   42,
			value: 43,
			want:  false,
		},
		{
			name:  "false - float64 / int",
			ref:   42.1,
			value: 42,
			want:  false,
		},
		{
			name:  "false - float64 / float64",
			ref:   42.0,
			value: 42.1,
			want:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &exact{}
			res, err := e.Evaluate(tt.ref, tt.value)
			require.NoError(t, err)
			require.Equal(t, tt.want, res, "Evaluate() = %v, want %v", tt.value, tt.want)
		})
	}
}

func TestNot_Evaluate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		internal Evaluator
		ref      interface{}
		value    interface{}
		want     bool
		wantErr  bool
	}{
		{
			name:     "true",
			internal: &exact{},
			ref:      1,
			value:    2,
			want:     true,
		},
		{
			name:     "false",
			internal: &exact{},
			ref:      1,
			value:    1,
			want:     false,
		},
		{
			name:     "error",
			internal: &evalRegexp{},
			ref:      "(",
			value:    "some value",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &not{
				Internal: tt.internal,
			}
			res, err := e.Evaluate(tt.ref, tt.value)

			if tt.wantErr {
				require.Error(t, err, "Apply() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, res, "Evaluate = %v, want %v", res, tt.want)
			}
		})
	}
}
