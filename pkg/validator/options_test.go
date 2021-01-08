package validator

import (
	"testing"

	vt "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestWithDefaultTranslations(t *testing.T) {
	tests := []struct {
		name string
		v    *Validator
	}{
		{
			name: "success",
			v: &Validator{
				V: vt.New(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := WithDefaultTranslations()(tt.v)
			require.NoError(t, err)
		})
	}
}

func TestWithFieldNameTag(t *testing.T) {
	tests := []struct {
		name string
		tag  string
	}{
		{
			name: "tag is empty string",
			tag:  "",
		},
		{
			name: "name with -",
			tag:  "-abc-efgh-ijkl",
		},
		{
			name: "success return name",
			tag:  "abcderfghijk",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := &Validator{
				V: vt.New(),
			}
			err := WithFieldNameTag(tt.tag)(v)
			require.NoError(t, err)
		})
	}
}
