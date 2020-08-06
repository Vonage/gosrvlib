package testutil

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"
)

// CaptureOutput hijacks and captures all log, stderr, stdout for testing
func CaptureOutput(fn func()) (string, error) {
	reader, writer, err := os.Pipe()
	if err != nil {
		return "", err
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)

	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		_, _ = io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	fn()
	_ = writer.Close()
	return <-out, nil
}
