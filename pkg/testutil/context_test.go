package testutil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nexmoinc/gosrvlib/pkg/httputil"
	"github.com/stretchr/testify/require"
)

func TestContextWithHTTPRouterParams(t *testing.T) {
	params := map[string]string{"test_arg_1": "test_val_1", "test_arg_2": "test_val_2"}

	r := httptest.NewRequest(http.MethodGet, "http://test.url.invalid", nil)
	r = r.WithContext(ContextWithHTTPRouterParams(context.Background(), params))

	require.Equal(t, "test_val_1", httputil.PathParam(r, "test_arg_1"))
	require.Equal(t, "test_val_2", httputil.PathParam(r, "test_arg_2"))
}
