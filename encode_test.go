package gocsv

import (
	"bytes"
	"encoding/csv"
	"testing"
)

type Sample struct {
	Foo string `csv:"foo"`
	Bar int    `csv:"BAR"`
	Baz string `csv:"Baz"`
}

func assertLine(t *testing.T, expected, actual []string) {
	if len(expected) != len(actual) {
		t.Fatalf("line length mismatch between expected: %d and actual: %d", len(expected), len(actual))
	}
	for i := range expected {
		if expected[i] != actual[i] {
			t.Fatalf("mismatch on field %d: %s != %s", i, expected[i], actual[i])
		}

	}
}

func Test_writeTo(t *testing.T) {
	b := bytes.Buffer{}
	w := csv.NewWriter(&b)
	s := []Sample{
		{Foo: "f", Bar: 1, Baz: "baz"},
		{Foo: "e", Bar: 3, Baz: "b"},
	}
	if err := writeTo(w, s); err != nil {
		t.Fatal(err)
	}

	lines, err := csv.NewReader(&b).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	assertLine(t, []string{"foo", "BAR", "Baz"}, lines[0])
	assertLine(t, []string{"f", "1", "baz"}, lines[1])
	assertLine(t, []string{"e", "3", "b"}, lines[2])
}
