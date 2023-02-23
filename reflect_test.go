package gocsv

import "testing"

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
