package gocsv

import (
	"strings"
	"testing"
)

type cyclicCSVRow struct {
	Name  string `csv:"name"`
	Inner *cyclicCSVInner
}

type cyclicCSVInner struct {
	Position string `csv:"position"`
	Parent   *cyclicCSVRow
}

func TestCyclicStructSchemaReturnsError(t *testing.T) {
	rows := []cyclicCSVRow{{Name: "one"}}
	if _, err := MarshalString(rows); err == nil || !strings.Contains(err.Error(), "cyclic struct type") {
		t.Fatalf("MarshalString should reject cyclic schema, got %v", err)
	}
	var decoded []cyclicCSVRow
	if err := UnmarshalString("name\none\n", &decoded); err == nil || !strings.Contains(err.Error(), "cyclic struct type") {
		t.Fatalf("UnmarshalString should reject cyclic schema, got %v", err)
	}
}

func TestExcludedCycleStillMarshalsNestedFields(t *testing.T) {
	type inner struct {
		Position string        `csv:"position"`
		Parent   *cyclicCSVRow `csv:"-"`
	}
	type row struct {
		Name  string `csv:"name"`
		Inner inner  `csv:"inner"`
	}
	got, err := MarshalString([]row{{Name: "one", Inner: inner{Position: "left"}}})
	if err != nil {
		t.Fatal(err)
	}
	if want := "name,inner.position\none,left\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func Test_fieldInfo_matchesKey(t *testing.T) {
	type fields struct {
		keys         []string
		omitEmpty    bool
		IndexChain   []int
		defaultValue string
		partial      bool
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "valid value",
			fields: fields{
				keys: []string{"date"},
			},
			args: args{"date"},
			want: true,
		},
		{
			name: "zero width space (U+200B)",
			fields: fields{
				keys: []string{"date"},
			},
			args: args{"\u200Bdate"},
			want: true,
		},
		{
			name: "zero width non-joiner (U+200C)",
			fields: fields{
				keys: []string{"date"},
			},
			args: args{"\u200Cdate"},
			want: true,
		},
		{
			name: "zero width joiner (U+200D)",
			fields: fields{
				keys: []string{"date"},
			},
			args: args{"\u200Ddate"},
			want: true,
		},
		{
			name: "zero width no-break space (U+FEFF)",
			fields: fields{
				keys: []string{"date"},
			},
			args: args{"\uFEFFdate"},
			want: true,
		},
		{
			name: "zero width no-break space (U+FEFF) in the middle of the string",
			fields: fields{
				keys: []string{"date"},
			},
			args: args{"da\uFEFFte"},
			want: true,
		},
		{
			name: "zero width no-break space (U+FEFF) in the end of the string",
			fields: fields{
				keys: []string{"date"},
			},
			args: args{"date\uFEFF"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fieldInfo{
				keys:         tt.fields.keys,
				omitEmpty:    tt.fields.omitEmpty,
				IndexChain:   tt.fields.IndexChain,
				defaultValue: tt.fields.defaultValue,
				partial:      tt.fields.partial,
			}
			if got := f.matchesKey(tt.args.key); got != tt.want {
				t.Errorf("matchesKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
