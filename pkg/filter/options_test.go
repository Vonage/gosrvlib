package filter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithFieldNameTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fieldTag string
		wantErr  bool
	}{
		{
			name:     "success",
			fieldTag: "json",
			wantErr:  false,
		},
		{
			name:     "error - empty field tag",
			fieldTag: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opt := WithFieldNameTag(tt.fieldTag)
			err := opt(&processor{})

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
