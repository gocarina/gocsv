package gocsv

import (
	"math"
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

func TestToIntLeadingZeros(t *testing.T) {
	testCases := []struct {
		field  string
		result int64
	}{
		{"08", 8},
		{"09", 9},
		{"010", 10},
		{"007", 7},
		{"0009", 9},
		{"-08", -8},
		{"-010", -10},
		{"0", 0},
		{"00", 0},
		{"010.9", 10},
	}

	for _, item := range testCases {
		out, err := toInt(item.field)
		if err != nil {
			t.Errorf("toInt(%q) returned an error: %s", item.field, err)
			continue
		}
		if out != item.result {
			t.Errorf("toInt(%q) = %d, want %d", item.field, out, item.result)
		}
	}
}

func TestToUintLeadingZeros(t *testing.T) {
	testCases := []struct {
		field  string
		result uint64
	}{
		{"08", 8},
		{"09", 9},
		{"010", 10},
		{"007", 7},
		{"0009", 9},
		{"0", 0},
		{"00", 0},
	}

	for _, item := range testCases {
		out, err := toUint(item.field)
		if err != nil {
			t.Errorf("toUint(%q) returned an error: %s", item.field, err)
			continue
		}
		if out != item.result {
			t.Errorf("toUint(%q) = %d, want %d", item.field, out, item.result)
		}
	}
}

func TestFloat32Extremes(t *testing.T) {
	type row struct {
		F float32 `csv:"f"`
	}

	// strconv's shortest representation of math.MaxFloat32 reads back as the
	// same float32 when parsed at 32 bits, but is above MaxFloat32 when parsed
	// at 64 bits; it must be accepted, as encoding/json and csvutil do.
	rows := []row{}
	if err := UnmarshalString("f\n3.4028235E+38\n-3.4028235E+38\n1E-45\n340282350000000000000000000000000000000\n", &rows); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []float32{math.MaxFloat32, -math.MaxFloat32, math.SmallestNonzeroFloat32, math.MaxFloat32}
	if len(rows) != len(want) {
		t.Fatalf("expected %d rows, got %d", len(want), len(rows))
	}
	for i, w := range want {
		if rows[i].F != w {
			t.Errorf("row %d: expected %v, got %v", i, w, rows[i].F)
		}
	}

	// values that do not fit in a float32 are still rejected
	var bad []row
	if err := UnmarshalString("f\n3.4028236E+38\n", &bad); err == nil {
		t.Errorf("expected a range error for 3.4028236E+38, got %v", bad)
	}

	// the library's own output round-trips
	out, err := MarshalString(&[]row{{F: math.MaxFloat32}, {F: -math.MaxFloat32}})
	if err != nil {
		t.Fatal(err)
	}
	var back []row
	if err := UnmarshalString(out, &back); err != nil {
		t.Fatalf("could not read own output %q: %v", out, err)
	}
	if len(back) != 2 || back[0].F != math.MaxFloat32 || back[1].F != -math.MaxFloat32 {
		t.Errorf("round trip changed the values: %v", back)
	}
}
