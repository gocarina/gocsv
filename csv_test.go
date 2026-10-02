package gocsv

import (
	"errors"
	"strings"
	"sync"
	"testing"
)

func TestUnmarshalToCallback_ReaderError(t *testing.T) {
	type Dummy struct{}
	var reader = &errorReader{}

	err := UnmarshalToCallback(reader, func(Dummy) {})
	if !errors.Is(err, readerErr) {
		t.Error("UnmarshalToCallback should return first reader error")
	}

	err = UnmarshalDecoderToCallback(newSimpleDecoderFromReader(reader), func(Dummy) {})
	if !errors.Is(err, readerErr) {
		t.Error("UnmarshalDecoderToCallback should return first reader error")
	}

	err = UnmarshalToCallbackWithError(reader, func(Dummy) error { return nil })
	if !errors.Is(err, readerErr) {
		t.Error("UnmarshalToCallbackWithError should return first reader error")
	}
}

func TestSetHeaderNormalizerWhileUnmarshaling(t *testing.T) {
	defer SetHeaderNormalizer(DefaultNameNormalizer())

	type row struct {
		ID   string `csv:"id"`
		Name string `csv:"name"`
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			SetHeaderNormalizer(strings.ToLower)
		}()
		go func() {
			defer wg.Done()
			var rows []row
			if err := UnmarshalString("id,name\n1,foo\n", &rows); err != nil {
				t.Error(err)
				return
			}
			if len(rows) != 1 || rows[0].ID != "1" || rows[0].Name != "foo" {
				t.Errorf("unexpected rows: %+v", rows)
			}
		}()
	}
	wg.Wait()
}

type errorReader struct{}

func (e *errorReader) Read([]byte) (n int, err error) {
	return 0, readerErr
}

var readerErr = errors.New("reader error")
