package gocsv

import (
	"encoding/csv"
	"io"
	"strings"
	"testing"
)

func TestUnmarshallerLongRow(t *testing.T) {
	type sample struct {
		FieldA string `csv:"field_a"`
		FieldB string `csv:"field_b"`
	}
	const csvContents = `field_a,field_b
a,b
c,d,e
`

	reader := csv.NewReader(strings.NewReader(csvContents))
	reader.FieldsPerRecord = -1
	um, err := NewUnmarshaller(reader, sample{})
	if err != nil {
		t.Fatalf("Error calling NewUnmarshaller: %#v", err)
	}

	obj, err := um.Read()
	if err != nil {
		t.Fatalf("Error calling Read(): %#v", err)
	}
	if obj.(sample).FieldA != "a" || obj.(sample).FieldB != "b" {
		t.Fatalf("Unepxected result from Read(): %#v", obj)
	}

	obj, err = um.Read()
	if err != nil {
		t.Fatalf("Error calling Read(): %#v", err)
	}
	if obj.(sample).FieldA != "c" || obj.(sample).FieldB != "d" {
		t.Fatalf("Unepxected result from Read(): %#v", obj)
	}

	obj, err = um.Read()
	if err != io.EOF {
		t.Fatalf("Unepxected result from Read(): (%#v, %#v)", obj, err)
	}
}
