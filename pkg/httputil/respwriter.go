package httputil

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

// NewWrapResponseWriter wraps an http.ResponseWriter with an enhanced proxy.
func NewWrapResponseWriter(w http.ResponseWriter) WrapResponseWriter {
	return &wrapResponseWriter{ResponseWriter: w}
}

// WrapResponseWriter is the interface defining the extendend functions of the proxy.
type WrapResponseWriter interface {
	http.ResponseWriter

	// StatusCode returns the HTTP status of the request.
	StatusCode() int

	// BytesCount returns the total number of bytes sent to the client.
	BytesCount() int

	// Tee sets a writer that will contain a copy of the bytes written to the response writer.
	Tee(io.Writer)
}

type wrapResponseWriter struct {
	http.ResponseWriter
	bytesCount    int
	headerWritten bool
	statusCode    int
	tee           io.Writer
}

func (b *wrapResponseWriter) BytesCount() int {
	return b.bytesCount
}

func (b *wrapResponseWriter) Flush() {
	b.headerWritten = true
	fl := b.ResponseWriter.(http.Flusher)
	fl.Flush()
}

// nolint:wrapcheck
func (b *wrapResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := b.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

// nolint:wrapcheck
func (b *wrapResponseWriter) Push(target string, opts *http.PushOptions) error {
	return b.ResponseWriter.(http.Pusher).Push(target, opts)
}

// nolint:wrapcheck
func (b *wrapResponseWriter) ReadFrom(r io.Reader) (int64, error) {
	if b.tee != nil {
		n, err := io.Copy(b, r)
		b.bytesCount += int(n)

		return n, err
	}

	rf := b.ResponseWriter.(io.ReaderFrom)

	b.maybeWriteHeader()

	n, err := rf.ReadFrom(r)

	b.bytesCount += int(n)

	return n, err
}

func (b *wrapResponseWriter) StatusCode() int {
	return b.statusCode
}

func (b *wrapResponseWriter) Tee(w io.Writer) {
	b.tee = w
}

func (b *wrapResponseWriter) Write(buf []byte) (int, error) {
	b.maybeWriteHeader()
	n, err := b.ResponseWriter.Write(buf)

	if b.tee != nil {
		_, teeErr := b.tee.Write(buf[:n])

		if err == nil {
			err = teeErr
		}
	}

	b.bytesCount += n

	return n, err
}

func (b *wrapResponseWriter) WriteHeader(code int) {
	if !b.headerWritten {
		b.statusCode = code
		b.headerWritten = true
		b.ResponseWriter.WriteHeader(code)
	}
}

func (b *wrapResponseWriter) maybeWriteHeader() {
	if !b.headerWritten {
		b.WriteHeader(http.StatusOK)
	}
}
