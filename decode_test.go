package gocsv

import "bytes"
import "testing"

func Test_readTo(t *testing.T) {
	b := bytes.NewBufferString(`foo,BAR,Baz
f,1,baz
e,3,b`)
	d := &decoder{in: b}

	var samples []Sample
	if err := d.readTo(&samples); err != nil {
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
