package gocsv

import (
	"strings"
	"testing"
)

// A struct field of type error makes decoding panic. The panic must reach the
// caller's goroutine, like it does for Unmarshal, so the caller can recover it.
type callbackPanicRecord struct {
	Err error `csv:"err"`
}

func TestCallbackPanicCanBeRecovered(t *testing.T) {
	const in = "err\nx\n"
	tests := map[string]func(){
		"UnmarshalToCallback": func() {
			_ = UnmarshalStringToCallback(in, func(callbackPanicRecord) {})
		},
		"UnmarshalToCallbackWithError": func() {
			_ = UnmarshalStringToCallbackWithError(in, func(callbackPanicRecord) error { return nil })
		},
		"UnmarshalDecoderToCallback": func() {
			_ = UnmarshalDecoderToCallback(newSimpleDecoderFromReader(strings.NewReader(in)), func(callbackPanicRecord) {})
		},
	}
	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected the decoding panic to reach the caller")
				}
			}()
			run()
			t.Fatal("expected a panic")
		})
	}
}
