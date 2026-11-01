package gocsv_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gocarina/gocsv"
)

func TestUnmarshalNestedParentAlias(t *testing.T) {
	type child struct {
		Value string `csv:"value"`
	}
	type row struct {
		Child child `csv:"primary,legacy"`
	}
	var rows []row
	err := gocsv.UnmarshalString("legacy.value\nkept\n", &rows)
	if err != nil || len(rows) != 1 || rows[0].Child.Value != "kept" {
		t.Fatalf("UnmarshalString: rows=%+v, err=%v; want Child.Value=kept", rows, err)
	}
}

func TestNestedAliasCombinations(t *testing.T) {
	type child struct {
		Value string `csv:"value,old_value"`
	}
	type row struct {
		Child child `csv:"primary,legacy,old"`
	}
	type pointerRow struct {
		Child *child `csv:"primary,legacy,old"`
	}
	type middle struct {
		Child *child `csv:"child,old_child"`
	}
	type deepRow struct {
		Middle *middle `csv:"primary,legacy,old"`
	}
	for _, tc := range []struct {
		name    string
		want    interface{}
		middles []string
		first   string
	}{
		{"value", []row{{Child: child{Value: "kept"}}}, []string{""}, "primary.value"},
		{"pointer", []pointerRow{{Child: &child{Value: "kept"}}}, []string{""}, "primary.value"},
		{"deep pointers", []deepRow{{Middle: &middle{Child: &child{Value: "kept"}}}}, []string{"child.", "old_child."}, "primary.child.value"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, parent := range []string{"primary", "legacy", "old"} {
				for _, mid := range tc.middles {
					for _, leaf := range []string{"value", "old_value"} {
						header := parent + "." + mid + leaf
						t.Run(header, func(t *testing.T) {
							out := reflect.New(reflect.TypeOf(tc.want))
							if err := gocsv.UnmarshalString(header+"\nkept\n", out.Interface()); err != nil {
								t.Fatal(err)
							}
							if got := out.Elem().Interface(); !reflect.DeepEqual(got, tc.want) {
								t.Fatalf("decoded %+v, want %+v", got, tc.want)
							}
						})
					}
				}
			}
			got, err := gocsv.MarshalString(tc.want)
			if want := tc.first + "\nkept\n"; err != nil || got != want {
				t.Fatalf("MarshalString = %q, %v; want %q", got, err, want)
			}
		})
	}
}

func TestNestedAliasesRejectAccumulatedPrefixes(t *testing.T) {
	type child struct {
		Value string `csv:"value,old_value"`
	}
	type row struct {
		Child child `csv:"primary,legacy,old"`
	}
	previous := gocsv.FailIfUnmatchedStructTags
	gocsv.FailIfUnmatchedStructTags = true
	defer func() { gocsv.FailIfUnmatchedStructTags = previous }()
	for _, header := range []string{
		"legacy.primary.value",
		"legacy.primary.old_value",
		"old.primary.value",
		"old.legacy.primary.value",
	} {
		t.Run(header, func(t *testing.T) {
			var rows []row
			err := gocsv.UnmarshalString(header+"\nkept\n", &rows)
			if !errors.Is(err, gocsv.ErrUnmatchedStructTags) {
				t.Fatalf("UnmarshalString(%q): rows=%+v, err=%v; want unmatched tags error", header, rows, err)
			}
		})
	}
}
