package gocsv

import (
	"encoding/csv"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestTwoColumnMapsRecordWidth(t *testing.T) {
	originalReader := selfCSVReader
	defer SetCSVReader(originalReader)
	SetCSVReader(func(in io.Reader) CSVReader {
		r := csv.NewReader(in)
		r.FieldsPerRecord = -1
		return r
	})

	parsers := []struct {
		name  string
		parse func(string) (map[string]string, error)
	}{
		{"CSVToMap", func(input string) (map[string]string, error) {
			return CSVToMap(strings.NewReader(input))
		}},
		{"UnmarshalCSVToMap", func(input string) (map[string]string, error) {
			var out map[string]string
			err := UnmarshalCSVToMap(getCSVReader(strings.NewReader(input)), &out)
			return out, err
		}},
	}
	for _, parser := range parsers {
		for _, tc := range []struct {
			name    string
			input   string
			want    map[string]string
			wantErr bool
		}{
			{"valid", "key,value\none,\"two,three\"\nfour,\n", map[string]string{"one": "two,three", "four": ""}, false},
			{"short", "key,value\none,two\nthree\n", nil, true},
			{"extra", "key,value\none,two\nthree,four,five\n", nil, true},
			{"header_only", "key,value\n", map[string]string{}, false},
		} {
			t.Run(parser.name+"/"+tc.name, func(t *testing.T) {
				got, err := parser.parse(tc.input)
				if (err != nil) != tc.wantErr {
					t.Fatalf("got error %v, want error %v", err, tc.wantErr)
				}
				if !tc.wantErr && !reflect.DeepEqual(got, tc.want) {
					t.Fatalf("got %#v, want %#v", got, tc.want)
				}
			})
		}
	}
}
