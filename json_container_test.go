package gocsv

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

type jsonTreeCSVRow struct {
	Name     string           `csv:"name" json:"name"`
	Children []jsonTreeCSVRow `csv:"children" json:"children"`
}

type jsonArrayCSVRow struct {
	Children [1]jsonArrayCSVChild `csv:"children"`
}

type jsonArrayCSVChild struct {
	Name   string           `json:"name"`
	Parent *jsonArrayCSVRow `json:"parent"`
}

func TestRecursiveJSONContainers(t *testing.T) {
	tests := []struct {
		name    string
		rows    interface{}
		records [][]string
	}{
		{
			name: "slice",
			rows: []jsonTreeCSVRow{{Name: "root", Children: []jsonTreeCSVRow{{Name: "leaf"}}}},
			records: [][]string{
				{"name", "children"},
				{"root", `[{"name":"leaf","children":null}]`},
			},
		},
		{
			name: "array",
			rows: []jsonArrayCSVRow{{Children: [1]jsonArrayCSVChild{{Name: "leaf"}}}},
			records: [][]string{
				{"children"},
				{`[{"name":"leaf","parent":null}]`},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var input bytes.Buffer
			if err := csv.NewWriter(&input).WriteAll(test.records); err != nil {
				t.Fatal(err)
			}
			t.Run("marshal", func(t *testing.T) {
				got, err := MarshalString(test.rows)
				if err != nil {
					t.Fatal(err)
				}
				if got != input.String() {
					t.Fatalf("got %q, want %q", got, input.String())
				}
			})
			t.Run("unmarshal", func(t *testing.T) {
				got := reflect.New(reflect.TypeOf(test.rows))
				if err := UnmarshalString(input.String(), got.Interface()); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got.Elem().Interface(), test.rows) {
					t.Fatalf("got %#v, want %#v", got.Elem().Interface(), test.rows)
				}
			})
		})
	}
}

func TestRecursiveJSONContainerRejectsCyclicValue(t *testing.T) {
	rows := []jsonTreeCSVRow{{Name: "root"}}
	rows[0].Children = rows
	_, err := MarshalString(rows)
	var jsonError *json.UnsupportedValueError
	if !errors.As(err, &jsonError) {
		t.Fatalf("got %v, want a JSON cyclic-value error", err)
	}
}

type indexedTreeCSVRow struct {
	Name     string              `csv:"name"`
	Children []indexedTreeCSVRow `csv:"children" csv[]:"1"`
}

func TestIndexedRecursiveContainerStillRejected(t *testing.T) {
	rows := []indexedTreeCSVRow{{Name: "root"}}
	if _, err := MarshalString(rows); err == nil || !strings.Contains(err.Error(), "cyclic struct type") {
		t.Fatalf("MarshalString should reject cyclic CSV expansion, got %v", err)
	}
	var decoded []indexedTreeCSVRow
	if err := UnmarshalString("name\nroot\n", &decoded); err == nil || !strings.Contains(err.Error(), "cyclic struct type") {
		t.Fatalf("UnmarshalString should reject cyclic CSV expansion, got %v", err)
	}
}
