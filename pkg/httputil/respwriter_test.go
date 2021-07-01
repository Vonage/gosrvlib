package httputil

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func newMockResponseWriter() *mockResponseWriter {
	buf := bytes.NewBuffer([]byte{})
	return &mockResponseWriter{Buffer: buf}
}

type mockResponseWriter struct {
	*bytes.Buffer
	hijackCalled bool
	pushCalled   bool
}

func (rw *mockResponseWriter) Header() http.Header {
	return nil
}

// nolint:wrapcheck
func (rw *mockResponseWriter) Write(in []byte) (int, error) {
	return rw.Buffer.Write(in)
}

func (rw *mockResponseWriter) WriteHeader(statusCode int) {

}

func (rw *mockResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	rw.hijackCalled = true
	return nil, nil, nil
}

func (rw *mockResponseWriter) Push(target string, opts *http.PushOptions) error {
	rw.pushCalled = true
	return nil
}

func TestNewWrapResponseWriter(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	ww := NewWrapResponseWriter(rr)
	require.NotNil(t, ww)
	require.Equal(t, reflect.ValueOf(rr).Pointer(), reflect.ValueOf(ww.(*wrapResponseWriter).ResponseWriter).Pointer())
}

func Test_wrapResponseWriter_BytesCount(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	ww := wrapResponseWriter{ResponseWriter: rr}
	count, err := ww.Write([]byte("test-counter"))
	require.Equal(t, 12, count)
	require.NoError(t, err)
	require.Equal(t, 12, ww.BytesCount())
}

func Test_wrapResponseWriter_Flush(t *testing.T) {
	t.Parallel()

	ww := wrapResponseWriter{ResponseWriter: httptest.NewRecorder()}
	ww.Flush()
	require.True(t, ww.headerWritten, "expected flush to set headerWritten=true")
}

func Test_wrapResponseWriter_StatusCode(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	ww := wrapResponseWriter{ResponseWriter: rr}
	ww.WriteHeader(97)
	require.Equal(t, 97, ww.StatusCode())
}

func Test_wrapResponseWriter_Tee(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	ww := wrapResponseWriter{ResponseWriter: rr}

	buf := bytes.NewBuffer([]byte{})
	ww.Tee(buf)

	count, err := ww.Write([]byte("tee"))
	require.Equal(t, 3, count)
	require.NoError(t, err)
	require.Equal(t, 3, ww.BytesCount())
	require.Equal(t, 3, buf.Len())
	require.Equal(t, "tee", buf.String())
}

func Test_wrapResponseWriter_Write(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	ww := wrapResponseWriter{ResponseWriter: rr}
	_, err := ww.Write([]byte("written"))
	require.NoError(t, err)
	require.Equal(t, 7, ww.BytesCount())
}

func Test_wrapResponseWriter_WriteHeader(t *testing.T) {
	t.Parallel()

	ww := wrapResponseWriter{ResponseWriter: httptest.NewRecorder()}
	ww.WriteHeader(19)
	require.Equal(t, 19, ww.StatusCode())
	ww.WriteHeader(41)
	require.Equal(t, 19, ww.StatusCode())
}

func Test_wrapResponseWriter_Hijack(t *testing.T) {
	t.Parallel()

	mock := newMockResponseWriter()
	ww := NewWrapResponseWriter(mock)
	require.NotNil(t, ww)

	_, _, err := ww.(*wrapResponseWriter).Hijack()
	require.NoError(t, err)
	require.True(t, mock.hijackCalled)
}

func Test_wrapResponseWriter_Push(t *testing.T) {
	t.Parallel()

	mock := newMockResponseWriter()
	ww := NewWrapResponseWriter(mock)
	require.NotNil(t, ww)

	_ = ww.(*wrapResponseWriter).Push("", &http.PushOptions{})

	require.True(t, mock.pushCalled)
}

func Test_wrapResponseWriter_ReadFrom(t *testing.T) {
	t.Parallel()

	// without tee
	mock := newMockResponseWriter()
	ww := NewWrapResponseWriter(mock)
	require.NotNil(t, ww)

	inputBuf := bytes.NewBufferString("0123456789")
	count, err := ww.(*wrapResponseWriter).ReadFrom(inputBuf)
	require.NoError(t, err)
	require.Equal(t, int64(10), count)

	// with tee writer
	mockTee := newMockResponseWriter()
	wwTee := NewWrapResponseWriter(mockTee)
	require.NotNil(t, wwTee)

	teeBuf := bytes.NewBuffer([]byte{})
	wwTee.Tee(teeBuf)

	inputBufTee := bytes.NewBufferString("0123456789")
	countTee, err := wwTee.(*wrapResponseWriter).ReadFrom(inputBufTee)
	require.NoError(t, err)
	require.Equal(t, int64(10), countTee)
	require.Equal(t, "0123456789", teeBuf.String())
	require.True(t, wwTee.(*wrapResponseWriter).headerWritten)
}
