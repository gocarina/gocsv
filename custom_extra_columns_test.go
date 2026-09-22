package gocsv

import (
	"encoding/csv"
	"reflect"
	"strings"
	"testing"
)

type extraColumnsCustomRecord struct {
	Name string   `csv:"name"`
	Seen []string `csv:"-"`
}

func (r *extraColumnsCustomRecord) UnmarshalCSVWithFields(key, value string) error {
	r.Seen = append(r.Seen, key+"="+value)
	if key == "name" {
		r.Name = value
	}
	return nil
}

func TestCustomUnmarshalIgnoresColumnsWithoutHeaders(t *testing.T) {
	for _, streaming := range []bool{false, true} {
		t.Run(map[bool]string{false: "slice", true: "channel"}[streaming], func(t *testing.T) {
			reader := csv.NewReader(strings.NewReader("name,other\nAlice,x,extra\nBob,y\n"))
			reader.FieldsPerRecord = -1
			decoder := NewSimpleDecoderFromCSVReader(reader)
			var records []extraColumnsCustomRecord
			if streaming {
				c := make(chan extraColumnsCustomRecord, 2)
				if err := UnmarshalDecoderToChan(decoder, c); err != nil {
					t.Fatal(err)
				}
				for r := range c {
					records = append(records, r)
				}
			} else if err := UnmarshalDecoder(decoder, &records); err != nil {
				t.Fatal(err)
			}
			want := []extraColumnsCustomRecord{{"Alice", []string{"name=Alice", "other=x"}}, {"Bob", []string{"name=Bob", "other=y"}}}
			if !reflect.DeepEqual(records, want) {
				t.Fatalf("got %#v, want %#v", records, want)
			}
		})
	}
}

func TestReadUnmatchedIgnoresColumnsWithoutHeaders(t *testing.T) {
	reader := csv.NewReader(strings.NewReader("name,other\nAlice,x,extra\n"))
	reader.FieldsPerRecord = -1
	type record struct {
		Name string `csv:"name"`
	}
	um, err := NewUnmarshaller(reader, record{})
	if err != nil {
		t.Fatal(err)
	}
	value, unmatched, err := um.ReadUnmatched()
	if err != nil {
		t.Fatal(err)
	}
	if value.(record).Name != "Alice" || !reflect.DeepEqual(unmatched, map[string]string{"other": "x"}) {
		t.Fatalf("got %#v, %#v", value, unmatched)
	}
}
