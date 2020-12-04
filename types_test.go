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

func (s sampleTypeUnmarshaller) MarshalCSV() (string, error) {
	return s.val, nil
}

type sampleTextUnmarshaller struct {
	val []byte
}

func (s *sampleTextUnmarshaller) UnmarshalText(text []byte) error {
	s.val = text
	return nil
}

func (s sampleTextUnmarshaller) MarshalText() ([]byte, error) {
	return s.val, nil
}

type sampleStringer string

func (s sampleStringer) String() string {
	return string(s)
}

type stringAlias string
type customStringAlias string

func (s customStringAlias) MarshalCSV() (string, error) {
	return `"` + string(s) + `"`, nil
}

func Test_getFieldAsString_CustomStringAlias(t *testing.T) {
	s, err := getFieldAsString(reflect.ValueOf(customStringAlias("foo")))
	if err != nil {
		t.Fatalf("getFieldAsString failure: %s", err)
	}

	if string(s) != `"foo"` {
		t.Fatalf(`expected "foo" got %s`, s)
	}

	s, err = getFieldAsString(reflect.ValueOf(stringAlias("foo")))
	if err != nil {
		t.Fatalf("getFieldAsString failure: %s", err)
	}

	if string(s) != `foo` {
		t.Fatalf(`expected foo got %s`, s)
	}
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

func Benchmark_marshall_TypeMarshaller(b *testing.B) {
	sample := sampleTypeUnmarshaller{"foo"}
	val := reflect.ValueOf(&sample)
	for n := 0; n < b.N; n++ {
		_, err := marshall(val)
		if err != nil {
			b.Fatalf("marshall error: %s", err.Error())
		}
	}
}

func Benchmark_marshall_TextMarshaller(b *testing.B) {
	sample := sampleTextUnmarshaller{[]byte("foo")}
	val := reflect.ValueOf(&sample)
	for n := 0; n < b.N; n++ {
		_, err := marshall(val)
		if err != nil {
			b.Fatalf("marshall error: %s", err.Error())
		}
	}
}

func Benchmark_marshall_Stringer(b *testing.B) {
	sample := sampleStringer("foo")
	val := reflect.ValueOf(&sample)
	for n := 0; n < b.N; n++ {
		_, err := marshall(val)
		if err != nil {
			b.Fatalf("marshall error: %s", err.Error())
		}
	}
}

func TestToInt(t *testing.T) {
	TestCase := []struct {
		field  string
		result int
		err    error
	}{
		{"123.2", 123, nil},
		{"123", 123, nil},
		{"1.2.3", 1, nil},
		{"0.123", 0, nil},
	}

	for idx, item := range TestCase {
		out, err := toInt(item.field)
		if err != nil {
			t.Fatal(err)
		}
		if int(out) != item.result {
			t.Fatal(idx, "result not equal")
		}
	}
}
