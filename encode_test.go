package gocsv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"
)

func assertLine(t *testing.T, expected, actual []string) {
	if len(expected) != len(actual) {
		t.Fatalf("line length mismatch between expected: %d and actual: %d\nExpected:\n%v\nActual:\n%v\n", len(expected), len(actual), expected, actual)
	}
	for i := range expected {
		if expected[i] != actual[i] {
			t.Fatalf("mismatch on field %d at line `%s`: %s != %s", i, expected, expected[i], actual[i])
		}
	}
}

func Test_writeTo(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}
	blah := 2
	sptr := "*string"
	s := []Sample{
		{Foo: "f", Bar: 1, Baz: "baz", Frop: 0.1, Blah: &blah, SPtr: &sptr},
		{Foo: "e", Bar: 3, Baz: "b", Frop: 6.0 / 13, Blah: nil, SPtr: nil},
	}
	if err := writeTo(NewSafeCSVWriter(csv.NewWriter(e.out)), s, false); err != nil {
		t.Fatal(err)
	}

	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	assertLine(t, []string{"foo", "BAR", "Baz", "Quux", "Blah", "SPtr", "Omit"}, lines[0])
	assertLine(t, []string{"f", "1", "baz", "0.1", "2", "*string", ""}, lines[1])
	assertLine(t, []string{"e", "3", "b", "0.46153846153846156", "", "", ""}, lines[2])
}

func Test_writeTo_Time(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}
	d := time.Unix(60, 0)
	s := []DateTime{
		{Foo: d},
	}
	if err := writeTo(NewSafeCSVWriter(csv.NewWriter(e.out)), s, true); err != nil {
		t.Fatal(err)
	}

	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	ft := time.Now()
	err = ft.UnmarshalText([]byte(lines[0][0]))
	if err != nil {
		t.Fatal(err)
	}
	if ft.Sub(d) != 0 {
		t.Fatalf("Dates doesn't match: %s and actual: %s", d, d)
	}

	m, _ := d.MarshalText()
	assertLine(t, []string{string(m)}, lines[0])
}

func Test_writeTo_NoHeaders(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}
	blah := 2
	sptr := "*string"
	s := []Sample{
		{Foo: "f", Bar: 1, Baz: "baz", Frop: 0.1, Blah: &blah, SPtr: &sptr},
		{Foo: "e", Bar: 3, Baz: "b", Frop: 6.0 / 13, Blah: nil, SPtr: nil},
	}
	if err := writeTo(NewSafeCSVWriter(csv.NewWriter(e.out)), s, true); err != nil {
		t.Fatal(err)
	}

	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	assertLine(t, []string{"f", "1", "baz", "0.1", "2", "*string", ""}, lines[0])
	assertLine(t, []string{"e", "3", "b", "0.46153846153846156", "", "", ""}, lines[1])
}

func Test_writeTo_multipleTags(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}
	s := []MultiTagSample{
		{Foo: "abc", Bar: 123},
		{Foo: "def", Bar: 234},
	}
	if err := writeTo(NewSafeCSVWriter(csv.NewWriter(e.out)), s, false); err != nil {
		t.Fatal(err)
	}

	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	// the first tag for each field is the encoding CSV header
	assertLine(t, []string{"Baz", "BAR"}, lines[0])
	assertLine(t, []string{"abc", "123"}, lines[1])
	assertLine(t, []string{"def", "234"}, lines[2])
}

func Test_writeTo_slice(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}

	type TestType struct {
		Key   string
		Items []int
	}

	s := []TestType{
		{
			Key:   "test1",
			Items: []int{1, 2, 3},
		},
		{
			Key:   "test2",
			Items: []int{4, 5, 6},
		},
	}

	if err := writeTo(NewSafeCSVWriter(csv.NewWriter(e.out)), s, false); err != nil {
		t.Fatal(err)
	}

	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}

	assertLine(t, []string{"Key", "Items"}, lines[0])
	assertLine(t, []string{"test1", "[1,2,3]"}, lines[1])
	assertLine(t, []string{"test2", "[4,5,6]"}, lines[2])
}

func Test_writeTo_slice_structs(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}
	s := []SliceStructSample{
		{
			Slice: []SliceStruct{
				{String: "s1", Float: 1.1},
				{String: "s2", Float: 2.2},
				{String: "nope", Float: 3.3},
			},
			Slice2: []SliceStruct{
				{String: "s1", Float: 1.1},
				{String: "s2", Float: 2.2},
				{String: "nope", Float: 3.3},
			},
			SimpleSlice: []int{1, 2, 3, 4, 5},
			Array: [2]SliceStruct{
				{String: "s3", Float: 3.3},
				{String: "s4", Float: 4.4},
			},
		},
	}
	if err := writeTo(NewSafeCSVWriter(csv.NewWriter(e.out)), s, false); err != nil {
		t.Fatal(err)
	}

	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	assertLine(t, []string{"s[0].s", "s[0].f", "s[1].s", "s[1].f", "sliceText", "ints[0]", "ints[1]", "ints[2]", "a[0].s", "a[0].f", "a[1].s", "a[1].f"}, lines[0])
	assertLine(t, []string{"s1", "1.1", "s2", "2.2", "[{\"String\":\"s1\",\"Float\":1.1},{\"String\":\"s2\",\"Float\":2.2},{\"String\":\"nope\",\"Float\":3.3}]", "1", "2", "3", "s3", "3.3", "s4", "4.4"}, lines[1])
}

func Test_writeTo_embed(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}
	blah := 2
	sptr := "*string"
	s := []EmbedSample{
		{
			Qux:    "aaa",
			Sample: Sample{Foo: "f", Bar: 1, Baz: "baz", Frop: 0.2, Blah: &blah, SPtr: &sptr},
			Ignore: "shouldn't be marshalled",
			Quux:   "zzz",
			Grault: math.Pi,
		},
	}
	if err := writeTo(NewSafeCSVWriter(csv.NewWriter(e.out)), s, false); err != nil {
		t.Fatal(err)
	}

	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	assertLine(t, []string{"first", "foo", "BAR", "Baz", "Quux", "Blah", "SPtr", "Omit", "garply", "last"}, lines[0])
	assertLine(t, []string{"aaa", "f", "1", "baz", "0.2", "2", "*string", "", "3.141592653589793", "zzz"}, lines[1])
}

func Test_writeTo_embedptr(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}
	blah := 2
	sptr := "*string"
	s := []EmbedPtrSample{
		{
			Qux:    "aaa",
			Sample: &Sample{Foo: "f", Bar: 1, Baz: "baz", Frop: 0.2, Blah: &blah, SPtr: &sptr},
			Ignore: "shouldn't be marshalled",
			Quux:   "zzz",
			Grault: math.Pi,
		},
	}
	if err := writeTo(NewSafeCSVWriter(csv.NewWriter(e.out)), s, false); err != nil {
		t.Fatal(err)
	}

	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	assertLine(t, []string{"first", "foo", "BAR", "Baz", "Quux", "Blah", "SPtr", "Omit", "garply", "last"}, lines[0])
	assertLine(t, []string{"aaa", "f", "1", "baz", "0.2", "2", "*string", "", "3.141592653589793", "zzz"}, lines[1])
}

func Test_writeTo_embedptr_nil(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}
	s := []EmbedPtrSample{
		{},
	}
	if err := writeTo(NewSafeCSVWriter(csv.NewWriter(e.out)), s, false); err != nil {
		t.Fatal(err)
	}

	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	assertLine(t, []string{"first", "foo", "BAR", "Baz", "Quux", "Blah", "SPtr", "Omit", "garply", "last"}, lines[0])
	assertLine(t, []string{"", "", "", "", "", "", "", "", "0", ""}, lines[1])
}

func Test_writeTo_embedmarshal(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}
	s := []EmbedMarshal{
		{
			Foo: &MarshalSample{Dummy: "bar"},
		},
	}
	if err := writeTo(NewSafeCSVWriter(csv.NewWriter(e.out)), s, false); err != nil {
		t.Fatal(err)
	}

	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	assertLine(t, []string{"foo"}, lines[0])
	assertLine(t, []string{"bar"}, lines[1])

}

func Test_writeTo_embedmarshalCSV(t *testing.T) {

	// First, create our test data
	b := new(bytes.Buffer)
	e := &encoder{out: b}
	s := []*EmbedMarshalCSV{
		{
			Symbol: "test",
			Timestamp: &MarshalCSVSample{
				Seconds: 1656460798,
				Nanos:   693201614,
			},
		},
	}

	// Next, attempt to write our test data to a CSV format
	if err := writeTo(NewSafeCSVWriter(csv.NewWriter(e.out)), s, false); err != nil {
		t.Fatal(err)
	}

	// Now, read in the data we just wrote
	lines, err := csv.NewReader(b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	// Finally, verify the structure of the data
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	assertLine(t, []string{"symbol", "timestamp"}, lines[0])
	assertLine(t, []string{"test", "1656460798.693201614"}, lines[1])
}

func Test_writeTo_complex_embed(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}
	sptr := "*string"
	sfs := []SkipFieldSample{
		{
			EmbedSample: EmbedSample{
				Qux: "aaa",
				Sample: Sample{
					Foo:  "bbb",
					Bar:  111,
					Baz:  "ddd",
					Frop: 1.2e22,
					Blah: nil,
					SPtr: &sptr,
				},
				Ignore: "eee",
				Grault: 0.1,
				Quux:   "fff",
			},
			MoreIgnore: "ggg",
			Corge:      "hhh",
		},
	}
	if err := writeTo(NewSafeCSVWriter(csv.NewWriter(e.out)), sfs, false); err != nil {
		t.Fatal(err)
	}
	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	assertLine(t, []string{"first", "foo", "BAR", "Baz", "Quux", "Blah", "SPtr", "Omit", "garply", "last", "abc"}, lines[0])
	assertLine(t, []string{"aaa", "bbb", "111", "ddd", "12000000000000000000000", "", "*string", "", "0.1", "fff", "hhh"}, lines[1])
}

func Test_writeTo_complex_inner_struct_embed(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}
	sfs := []Level0Struct{
		{
			Level0Field: level1Struct{
				Level1Field: level2Struct{
					InnerStruct{
						BoolIgnoreField0: false,
						BoolField1:       false,
						StringField2:     "email1",
					},
				},
			},
		},
		{
			Level0Field: level1Struct{
				Level1Field: level2Struct{
					InnerStruct{
						BoolIgnoreField0: false,
						BoolField1:       true,
						StringField2:     "email2",
					},
				},
			},
		},
	}

	if err := writeTo(NewSafeCSVWriter(csv.NewWriter(e.out)), sfs, true); err != nil {
		t.Fatal(err)
	}
	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	assertLine(t, []string{"false", "email1"}, lines[0])
	assertLine(t, []string{"true", "email2"}, lines[1])
}

func Test_writeToChan(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}
	c := make(chan interface{})
	sptr := "*string"
	go func() {
		for i := 0; i < 100; i++ {
			v := Sample{Foo: "f", Bar: i, Baz: "baz" + strconv.Itoa(i), Frop: float64(i), Blah: nil, SPtr: &sptr}
			c <- v
		}
		close(c)
	}()
	if err := MarshalChan(c, NewSafeCSVWriter(csv.NewWriter(e.out))); err != nil {
		t.Fatal(err)
	}
	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 101 {
		t.Fatalf("expected 100 lines, got %d", len(lines))
	}
	for i, l := range lines {
		if i == 0 {
			assertLine(t, []string{"foo", "BAR", "Baz", "Quux", "Blah", "SPtr", "Omit"}, l)
			continue
		}
		assertLine(t, []string{"f", strconv.Itoa(i - 1), "baz" + strconv.Itoa(i-1), strconv.FormatFloat(float64(i-1), 'f', -1, 64), "", "*string", ""}, l)
	}
}

func Test_MarshalChan_ClosedChannel(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}
	c := make(chan interface{})
	close(c)

	if err := MarshalChan(c, NewSafeCSVWriter(csv.NewWriter(e.out))); !errors.Is(err, ErrChannelIsClosed) {
		t.Fatal(err)
	}
}

// TestRenamedTypes tests for marshaling functions on redefined basic types.
func TestRenamedTypesMarshal(t *testing.T) {
	samples := []RenamedSample{
		{RenamedFloatUnmarshaler: 1.4, RenamedFloatDefault: 1.5},
		{RenamedFloatUnmarshaler: 2.3, RenamedFloatDefault: 2.4},
	}

	SetCSVWriter(func(out io.Writer) *SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = ';'
		return NewSafeCSVWriter(writer)
	})
	// Switch back to default for tests executed after this
	defer SetCSVWriter(DefaultCSVWriter)

	csvContent, err := MarshalString(&samples)
	if err != nil {
		t.Fatal(err)
	}
	if csvContent != "foo;bar\n1,4;1.5\n2,3;2.4\n" {
		t.Fatalf("Error marshaling floats with , as separator. Expected \nfoo;bar\n1,4;1.5\n2,3;2.4\ngot:\n%v", csvContent)
	}

	// Test that errors raised by MarshalCSV are correctly reported
	samples = []RenamedSample{
		{RenamedFloatUnmarshaler: 4.2, RenamedFloatDefault: 1.5},
	}
	_, err = MarshalString(&samples)
	if _, ok := err.(MarshalError); !ok {
		t.Fatalf("Expected UnmarshalError, got %v", err)
	}
}

// TestCustomTagSeparatorMarshal tests for custom tag separator in marshalling.
func TestCustomTagSeparatorMarshal(t *testing.T) {
	samples := []RenamedSample{
		{RenamedFloatUnmarshaler: 1.4, RenamedFloatDefault: 1.5},
		{RenamedFloatUnmarshaler: 2.3, RenamedFloatDefault: 2.4},
	}

	TagSeparator = " | "
	// Switch back to default TagSeparator after this
	defer func() {
		TagSeparator = ","
	}()

	csvContent, err := MarshalString(&samples)
	if err != nil {
		t.Fatal(err)
	}
	if csvContent != "foo|bar\n1,4|1.5\n2,3|2.4\n" {
		t.Fatalf("Error marshaling floats with , as separator. Expected \nfoo|bar\n1,4|1.5\n2,3|2.4\ngot:\n%v", csvContent)
	}
}

func (rf *RenamedFloat64Unmarshaler) MarshalCSV() (csv string, err error) {
	if *rf == RenamedFloat64Unmarshaler(4.2) {
		return "", MarshalError{"Test error: Invalid float 4.2"}
	}
	csv = strconv.FormatFloat(float64(*rf), 'f', 1, 64)
	csv = strings.Replace(csv, ".", ",", -1)
	return csv, nil
}

type MarshalError struct {
	msg string
}

func (e MarshalError) Error() string {
	return e.msg
}

func Benchmark_MarshalCSVWithoutHeaders(b *testing.B) {
	dst := NewSafeCSVWriter(csv.NewWriter(ioutil.Discard))
	for n := 0; n < b.N; n++ {
		err := MarshalCSVWithoutHeaders([]Sample{{}}, dst)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Test_writeTo_nested_struct(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}
	s := []NestedSample{
		{
			Inner1: InnerStruct{
				BoolIgnoreField0: false,
				BoolField1:       false,
				StringField2:     "email_one",
			},
			Inner2: InnerStruct{
				BoolIgnoreField0: true,
				BoolField1:       true,
				StringField2:     "email_two",
			},
			InnerIgnore: InnerStruct{
				BoolIgnoreField0: true,
				BoolField1:       false,
				StringField2:     "email_ignore",
			},
			Inner3: NestedEmbedSample{InnerStruct{
				BoolIgnoreField0: true,
				BoolField1:       false,
				StringField2:     "email_three",
			}},
		},
	}
	if err := writeTo(NewSafeCSVWriter(csv.NewWriter(e.out)), s, false); err != nil {
		t.Fatal(err)
	}

	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	assertLine(t, []string{"one.boolField1", "one.stringField2", "two.boolField1", "two.stringField2", "three.boolField1", "three.stringField2"}, lines[0])
	assertLine(t, []string{"false", "email_one", "true", "email_two", "false", "email_three"}, lines[1])
	t.Run("struct is a child element ", func(t *testing.T) {
		type Nested struct {
			UpdatedAt *time.Time `csv:"updated_at"`
		}
		type Row struct {
			Name string `csv:"name"`
			Data Nested `csv:"data"`
		}
		data := time.Date(2022, 12, 11, 13, 14, 15, 123456, time.UTC)
		dat := []Row{{Name: "matthew", Data: Nested{UpdatedAt: &data}}}
		got, err := MarshalString(dat)
		if err != nil {
			t.Fatal(err)
		}
		want := `name,data.updated_at
matthew,2022-12-11T13:14:15.000123456Z
`
		if got != want {
			t.Fatalf("Error marshaling floats with , as separator. Expected \n%s\ngot:\n%s", want, got)
		}
	})

}

func Test_non_marshaling_nested_fields_are_prefixed(t *testing.T) {
	b := bytes.Buffer{}
	e := &encoder{out: &b}
	time1 := time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2022, 2, 19, 0, 0, 0, 0, time.UTC)
	time3 := time.Date(2023, 2, 19, 0, 0, 0, 0, time.UTC)

	s := []SameNameStruct{
		{
			Inner2: &InnerStruct2{
				Bar:    "bar1",
				Inner3: InnerStruct3{Bar: "bar2", Foo: time1},
			},
			Inner3: InnerStruct3{Bar: "bar3", Foo: time2},
			Foo:    time3,
		},
	}
	if err := writeTo(NewSafeCSVWriter(csv.NewWriter(e.out)), s, false); err != nil {
		t.Fatal(err)
	}
	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	// The headers should contain the struct path prefixes even if struct is responsible for its marshalling like time.Time
	assertLine(t, []string{"Inner2.Bar", "Inner2.Inner3.Bar", "Inner2.Inner3.Foo", "Inner3.Bar", "Inner3.Foo", "Foo"}, lines[0])
	assertLine(t, []string{"bar1", "bar2", "2021-02-19T00:00:00Z", "bar3", "2022-02-19T00:00:00Z", "2023-02-19T00:00:00Z"}, lines[1])
}
