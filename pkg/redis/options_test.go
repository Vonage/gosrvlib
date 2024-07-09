package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WithMessageEncodeFunc(t *testing.T) {
	t.Parallel()

	ret := "test_data_001"
	f := func(_ context.Context, _ any) (string, error) {
		return ret, nil
	}

	conf := &cfg{}
	WithMessageEncodeFunc(f)(conf)

	d, err := conf.messageEncodeFunc(context.TODO(), "")
	require.NoError(t, err)
	require.Equal(t, ret, d)
}

func Test_WithMessageDecodeFunc(t *testing.T) {
	t.Parallel()

	f := func(_ context.Context, _ string, _ any) error {
		return nil
	}

	conf := &cfg{}
	WithMessageDecodeFunc(f)(conf)
	require.NoError(t, conf.messageDecodeFunc(context.TODO(), "", ""))
}
