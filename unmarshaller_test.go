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

func TestUnmarshallerRenormalizeHeaders(t *testing.T) {
	type sample struct {
		FieldA string `csv:"field_a_map"`
		FieldB string `csv:"field_b_map"`
	}
	const csvContents = `field_a,field_b
a,b
c,d,e
`

	headerNormalizer := func(headers []string) []string {
		normalizedHeaders := make([]string, len(headers))
		for i, header := range headers {
			normalizedHeader := header
			switch header {
			case "field_a":
				normalizedHeader = "field_a_map"
			case "field_b":
				normalizedHeader = "field_b_map"
			}
			normalizedHeaders[i] = normalizedHeader
		}
		return normalizedHeaders
	}

	reader := csv.NewReader(strings.NewReader(csvContents))
	reader.FieldsPerRecord = -1
	um, err := NewUnmarshaller(reader, sample{})
	if err != nil {
		t.Fatalf("Error calling NewUnmarshaller: %#v", err)
	}

	err = um.RenormalizeHeaders(headerNormalizer)
	if err != nil {
		t.Fatalf("Error calling RenormalizeHeaders: %#v", err)
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
