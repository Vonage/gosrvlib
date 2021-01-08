package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidationError_Error(t *testing.T) {
	t.Parallel()
	want := "mock_error"
	e := &ValidationError{Err: "mock_error"}
	got := e.Error()
	require.Equal(t, want, got, "Error() = %v, want %v", got, want)
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		want    bool
		wantErr bool
	}{
		{
			name:    "empty option success",
			opts:    nil,
			want:    true,
			wantErr: false,
		},
		{
			name: "applied opts returns error",
			opts: []Option{
				WithDefaultTranslations(),
				WithFieldNameTag("test_tag"),
				WithValidationTranslated("test_tag", nil, nil, nil),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "success with opts applied",
			opts: []Option{
				WithDefaultTranslations(),
				WithFieldNameTag("test_tag"),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := New(tt.opts...)

			if tt.wantErr {
				require.Nil(t, got, "New() returned Validator should be nil")
				require.Error(t, err, "New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NotNil(t, got, "New() returned Validator should not be nil")
			require.NoError(t, err, "New() unexpected error = %v", err)
		})
	}
}

func TestValidator_ValidateStruct(t *testing.T) {
	type baseConfig struct {
		Base1 string `mapstructure:"base1" validate:"required,url"`
		Base2 int    `mapstructure:"base2" validate:"required"`
	}
	type appConfig struct {
		Base    baseConfig `mapstructure:"base" validate:"required"`
		Config1 bool       `mapstructure:"config1" validate:"required"`
		Config2 string     `mapstructure:"config2" validate:"required"`
	}

	tests := []struct {
		name    string
		obj     interface{}
		wantErr bool
	}{
		{
			name: "success struct validated",
			obj: appConfig{
				Base: baseConfig{
					Base1: "https://test.ipify.url.invalid",
					Base2: 1234,
				},
				Config1: true,
				Config2: "test_string",
			},
			wantErr: false,
		},
		{
			name: "failed in validation",
			obj: appConfig{
				Base: baseConfig{Base1: "not_url"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			opts := []Option{
				WithFieldNameTag("mapstructure"),
				WithDefaultTranslations(),
			}
			v, _ := New(opts...)
			err := v.ValidateStruct(tt.obj)
			require.Equal(t, tt.wantErr, err != nil, "ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}
