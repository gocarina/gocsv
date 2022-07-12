package gocsv

import (
	"bytes"
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

func TestUnmarshalListOfStructsAfterMarshal(t *testing.T) {

	type Additional struct {
		Value string
	}

	type Option struct {
		Additional []*Additional
		Key        string
	}

	inData := []*Option{
		{
			Key: "test",
		},
	}

	// First, marshal our test data to a CSV format
	buffer := new(bytes.Buffer)
	innerWriter := csv.NewWriter(buffer)
	innerWriter.Comma = '|'
	csvWriter := NewSafeCSVWriter(innerWriter)
	if err := MarshalCSV(inData, csvWriter); err != nil {
		t.Fatalf("Error marshalling data to CSV: %#v", err)
	}

	if string(buffer.Bytes()) != "Additional|Key\nnull|test\n" {
		t.Fatalf("Marshalled data had an unexpected form of %s", buffer.Bytes())
	}

	// Next, attempt to unmarshal our test data from a CSV format
	var outData []*Option
	innerReader := csv.NewReader(buffer)
	innerReader.Comma = '|'
	if err := UnmarshalCSV(innerReader, &outData); err != nil {
		t.Fatalf("Error unmarshalling data from CSV: %#v", err)
	}

	// Finally, verify the data
	if len(outData) != 1 {
		t.Fatalf("Data expected to have one entry, had %d entries", len(outData))
	} else if len(outData[0].Additional) != 0 {
		t.Fatalf("Data Additional field expected to be empty, had length of %d", len(outData[0].Additional))
	} else if outData[0].Key != "test" {
		t.Fatalf("Data Key field did not contain expected value, had %q", outData[0].Key)
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
