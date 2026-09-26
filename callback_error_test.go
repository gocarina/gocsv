package gocsv

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type callbackErrorRecord struct {
	Value string `csv:"value"`
}

type callbackEOFReader struct {
	io.Reader
	eof chan struct{}
}

func (r *callbackEOFReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if err == io.EOF {
		close(r.eof)
	}
	return n, err
}

func TestCallbackErrorDrainsParser(t *testing.T) {
	for _, tail := range []string{"", "second\nthird\n", "\"bad\n"} {
		t.Run(tail, func(t *testing.T) {
			expected := errors.New("stop callback")
			reader := &callbackEOFReader{strings.NewReader("value\nfirst\n" + tail), make(chan struct{})}
			calls := 0
			err := UnmarshalToCallback(reader, func(callbackErrorRecord) error { calls++; return expected })
			if err != expected {
				t.Fatalf("got %v, want original callback error", err)
			}
			if calls != 1 {
				t.Fatalf("callback called %d times", calls)
			}
			select {
			case <-reader.eof:
			default:
				t.Fatal("parser was abandoned before reaching EOF")
			}
		})
	}
}
