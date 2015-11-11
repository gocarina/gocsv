package gocsv

import "bytes"
import "testing"

func Test_readTo(t *testing.T) {
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
}

func Test_readTo_complex_embed(t *testing.T) {
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
