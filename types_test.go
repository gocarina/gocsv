package gocsv

import (
	"reflect"
	"testing"
)

type sampleTypeUnmarshaller struct {
	val string
}

func (s *sampleTypeUnmarshaller) UnmarshalCSV(val string) error {
	s.val = val
	return nil
}

type sampleTextUnmarshaller struct {
	val []byte
}

func (s *sampleTextUnmarshaller) UnmarshalText(text []byte) error {
	s.val = text
	return nil
}

func Benchmark_unmarshall_TypeUnmarshaller(b *testing.B) {
	sample := sampleTypeUnmarshaller{}
	val := reflect.ValueOf(&sample)
	for n := 0; n < b.N; n++ {
		if err := unmarshall(val, "foo"); err != nil {
			b.Fatalf("unmarshall error: %s", err.Error())
		}
	}
}

func Benchmark_unmarshall_TextUnmarshaller(b *testing.B) {
	sample := sampleTextUnmarshaller{}
	val := reflect.ValueOf(&sample)
	for n := 0; n < b.N; n++ {
		if err := unmarshall(val, "foo"); err != nil {
			b.Fatalf("unmarshall error: %s", err.Error())
		}
	}
}
