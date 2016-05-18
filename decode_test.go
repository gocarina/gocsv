package gocsv

import (
	"bytes"
	"encoding/csv"
	"testing"
)

func resetFailIfUnmatchedStructTags(b bool) {
	FailIfUnmatchedStructTags = b
}

func Test_readTo(t *testing.T) {
	defer resetFailIfUnmatchedStructTags(FailIfUnmatchedStructTags)
	FailIfUnmatchedStructTags = false

	b := bytes.NewBufferString(`foo,BAR,Baz
f,1,baz
e,3,b`)
	d := &decoder{in: b}

	var samples []Sample
	if err := readTo(d, &samples); err != nil {
		t.Fatal(err)
	}
	if len(samples) != 2 {
		t.Fatalf("expected 2 sample instances, got %d", len(samples))
	}
	expected := Sample{Foo: "f", Bar: 1, Baz: "baz"}
	if expected != samples[0] {
		t.Fatalf("expected first sample %v, got %v", expected, samples[0])
	}
	expected = Sample{Foo: "e", Bar: 3, Baz: "b"}
	if expected != samples[1] {
		t.Fatalf("expected second sample %v, got %v", expected, samples[1])
	}

	b = bytes.NewBufferString(`foo,BAR,Baz
f,1,baz
e,BAD_INPUT,b`)
	d = &decoder{b}
	samples = []Sample{}
	err := readTo(d, &samples)
	if err == nil {
		t.Fatalf("Expected error from bad input, got: %+v", samples)
	}
	switch actualErr := err.(type) {
	case *csv.ParseError:
		if actualErr.Line != 3 {
			t.Fatalf("Expected csv.ParseError on line 3, got: %d", actualErr.Line)
		}
		if actualErr.Column != 2 {
			t.Fatalf("Expected csv.ParseError in column 2, got: %d", actualErr.Column)
		}
	default:
		t.Fatalf("incorrect error type: %T", err)
	}
}

func Test_readTo_multipleTags(t *testing.T) {
	b := bytes.NewBufferString(`Foo,bar
abc,123
def,234`)
	d := &decoder{in: b}

	var samples []MultiTagSample
	if err := readTo(d, &samples); err != nil {
		t.Fatal(err)
	}
	if len(samples) != 2 {
		t.Fatalf("expected 2 sample instances, got %d", len(samples))
	}
	expected := MultiTagSample{Foo: "abc", Bar: 123}
	if expected != samples[0] {
		t.Fatalf("expected first sample %v, got %v", expected, samples[0])
	}
	expected = MultiTagSample{Foo: "def", Bar: 234}
	if expected != samples[1] {
		t.Fatalf("expected second sample %v, got %v", expected, samples[1])
	}
}

func Test_readTo_complex_embed(t *testing.T) {
	defer resetFailIfUnmatchedStructTags(FailIfUnmatchedStructTags)
	FailIfUnmatchedStructTags = false

	b := bytes.NewBufferString(`first,foo,BAR,Baz,last,abc
aa,bb,11,cc,dd,ee
ff,gg,22,hh,ii,jj`)
	d := &decoder{in: b}

	var samples []SkipFieldSample
	if err := readTo(d, &samples); err != nil {
		t.Fatal(err)
	}
	if len(samples) != 2 {
		t.Fatalf("expected 2 sample instances, got %d", len(samples))
	}
	expected := SkipFieldSample{
		EmbedSample: EmbedSample{
			Qux: "aa",
			Sample: Sample{
				Foo: "bb",
				Bar: 11,
				Baz: "cc",
			},
			Quux: "dd",
		},
		Corge: "ee",
	}
	if expected != samples[0] {
		t.Fatalf("expected first sample %v, got %v", expected, samples[0])
	}
	expected = SkipFieldSample{
		EmbedSample: EmbedSample{
			Qux: "ff",
			Sample: Sample{
				Foo: "gg",
				Bar: 22,
				Baz: "hh",
			},
			Quux: "ii",
		},
		Corge: "jj",
	}
	if expected != samples[1] {
		t.Fatalf("expected first sample %v, got %v", expected, samples[1])
	}
}

func Test_maybeMissingStructFields(t *testing.T) {
	structTags := []fieldInfo{
		{keys: []string{"foo"}},
		{keys: []string{"bar"}},
		{keys: []string{"baz"}},
	}
	badHeaders := []string{"hi", "mom", "bacon"}
	goodHeaders := []string{"foo", "bar", "baz"}

	// no tags to match, expect no error
	if err := maybeMissingStructFields([]fieldInfo{}, goodHeaders); err != nil {
		t.Fatal(err)
	}

	// bad headers, expect an error
	if err := maybeMissingStructFields(structTags, badHeaders); err == nil {
		t.Fatal("expected an error, but no error found")
	}

	// good headers, expect no error
	if err := maybeMissingStructFields(structTags, goodHeaders); err != nil {
		t.Fatal(err)
	}

	// extra headers, but all structtags match; expect no error
	moarHeaders := append(goodHeaders, "qux", "quux", "corge", "grault")
	if err := maybeMissingStructFields(structTags, moarHeaders); err != nil {
		t.Fatal(err)
	}

	// not all structTags match, but there's plenty o' headers; expect
	// error
	mismatchedHeaders := []string{"foo", "qux", "quux", "corgi"}
	if err := maybeMissingStructFields(structTags, mismatchedHeaders); err == nil {
		t.Fatal("expected an error, but no error found")
	}
}
