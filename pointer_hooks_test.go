package gocsv

import (
	"encoding/csv"
	"errors"
	"reflect"
	"strings"
	"testing"
)

var errPointerHook = errors.New("rejected by row hook")

type pointerHookRow struct {
	Value string `csv:"value"`
}

func (r *pointerHookRow) UnmarshalCSVWithFields(key, value string) error {
	if value == "invalid" {
		return errPointerHook
	}
	r.Value += key + "=" + strings.ToUpper(value) + ";"
	return nil
}

func TestPointerRowHooks(t *testing.T) {
	const input = "value,extra\none,x\ntwo,y\n"
	want := []pointerHookRow{{"value=ONE;extra=X;"}, {"value=TWO;extra=Y;"}}
	for _, tc := range []struct {
		name string
		run  func(t *testing.T) ([]pointerHookRow, error)
	}{
		{"values", func(t *testing.T) ([]pointerHookRow, error) {
			var rows []pointerHookRow
			err := UnmarshalString(input, &rows)
			return rows, err
		}},
		{"pointers", func(t *testing.T) ([]pointerHookRow, error) {
			var rows []*pointerHookRow
			err := UnmarshalString(input, &rows)
			var values []pointerHookRow
			for _, r := range rows {
				values = append(values, *r)
			}
			if len(rows) == 2 && rows[0] == rows[1] {
				t.Fatal("rows share storage")
			}
			return values, err
		}},
		{"array", func(t *testing.T) ([]pointerHookRow, error) {
			var rows [2]*pointerHookRow
			err := UnmarshalString(input, &rows)
			var values []pointerHookRow
			for _, r := range rows {
				if r != nil {
					values = append(values, *r)
				}
			}
			return values, err
		}},
		{"channel", func(t *testing.T) ([]pointerHookRow, error) {
			ch := make(chan *pointerHookRow, 2)
			err := UnmarshalStringToChan(input, ch)
			var rows []pointerHookRow
			for r := range ch {
				rows = append(rows, *r)
			}
			return rows, err
		}},
		{"callback", func(t *testing.T) ([]pointerHookRow, error) {
			var rows []pointerHookRow
			err := UnmarshalStringToCallback(input, func(r *pointerHookRow) { rows = append(rows, *r) })
			return rows, err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.run(t)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got %#v, want %#v", got, want)
			}
		})
	}
}

func TestPointerRowHookErrors(t *testing.T) {
	const input = "value\ninvalid\n"
	for _, tc := range []struct {
		name string
		run  func(t *testing.T) error
	}{
		{"slice", func(t *testing.T) error { var rows []*pointerHookRow; return UnmarshalString(input, &rows) }},
		{"channel", func(t *testing.T) error {
			ch := make(chan *pointerHookRow, 1)
			err := UnmarshalStringToChan(input, ch)
			if len(ch) != 0 {
				t.Fatal("emitted rejected row")
			}
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.run(t)
			var pe *csv.ParseError
			if !errors.As(err, &pe) || !errors.Is(err, errPointerHook) {
				t.Fatalf("expected hook parse error, got %v", err)
			}
			if pe.Line != 2 || pe.Column != 1 {
				t.Fatalf("wrong location: %v", pe)
			}
		})
	}
}
