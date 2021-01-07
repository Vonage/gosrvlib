package validator

import (
	"testing"

	ut "github.com/go-playground/universal-translator"
	vt "github.com/go-playground/validator/v10"
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
	type args struct {
		opts []Option
	}
	tests := []struct {
		name           string
		args           args
		want           bool
		wantErr        bool
		wantValidate   *vt.Validate
		wantTranslator ut.Translator
	}{
		{
			name:    "empty option success",
			args:    args{nil},
			want:    true,
			wantErr: false,
		},
		{
			name: "applied opts returns error",
			args: args{
				opts: []Option{
					WithDefaultTranslations(),
					WithFieldNameTag("test_tag"),
					WithValidationTranslated("test_tag", nil, nil, nil),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "success with opts applied",
			args: args{
				opts: []Option{
					WithDefaultTranslations(),
					WithFieldNameTag("test_tag"),
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.opts...)

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
