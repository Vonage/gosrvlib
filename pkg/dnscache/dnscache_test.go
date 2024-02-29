package dnscache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	got := New(nil, 3, 1*time.Second)
	require.NotNil(t, got)

	require.NotNil(t, got.resolver)
	require.NotNil(t, got.mux)

	require.Equal(t, 3, got.size)
	require.Equal(t, 1*time.Second, got.ttl)

	require.NotNil(t, got.cache)
	require.Empty(t, got.cache)
}
