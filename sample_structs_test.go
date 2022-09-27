package gocsv

import (
	"fmt"
	"time"
)

type Sample struct {
	Foo  string  `csv:"foo"`
	Bar  int     `csv:"BAR"`
	Baz  string  `csv:"Baz"`
	Frop float64 `csv:"Quux"`
	Blah *int    `csv:"Blah"`
	SPtr *string `csv:"SPtr"`
	Omit *string `csv:"Omit,omitempty"`
}

type SliceSample struct {
	Slice []int `csv:"Slice"`
}

type SliceStructSample struct {
	Slice       []SliceStruct  `csv:"s,slice" csv[]:"2"`
	SimpleSlice []int          `csv:"ints" csv[]:"3"`
	Array       [2]SliceStruct `csv:"a,array" csv[]:"2"`
}

type SliceStruct struct {
	String string  `csv:"s,string"`
	Float  float64 `csv:"f,float"`
}

type EmbedSample struct {
	Qux string `csv:"first"`
	Sample
	Ignore string  `csv:"-"`
	Grault float64 `csv:"garply"`
	Quux   string  `csv:"last"`
}

type MarshalSample struct {
	Dummy string
}

func (m MarshalSample) MarshalText() ([]byte, error) {
	return []byte(m.Dummy), nil
}
func (m *MarshalSample) UnmarshalText(text []byte) error {
	m.Dummy = string(text)
	return nil
}

type EmbedMarshal struct {
	Foo *MarshalSample `csv:"foo"`
}

type MarshalCSVSample struct {
	Seconds int64
	Nanos   int32
}

func (timestamp *MarshalCSVSample) MarshalCSV() (string, error) {
	if timestamp == nil {
		return "", nil
	}

	return fmt.Sprintf("%d.%09d", timestamp.Seconds, timestamp.Nanos), nil
}

type EmbedMarshalCSV struct {
	Symbol    string            `csv:"symbol"`
	Timestamp *MarshalCSVSample `csv:"timestamp"`
}

type UnmarshalCSVSample struct {
	Timestamp int64
	Nanos     int32
}

func (timestamp *UnmarshalCSVSample) UnmarshalCSV(s string) error {
	ret := UnmarshalCSVSample{}
	_, err := fmt.Sscanf(s, "%d.%09d", &ret.Timestamp, &ret.Nanos)
	*timestamp = ret
	return err
}

type EmbedUnmarshalCSVWithClashingField struct {
	Symbol string

	// Clashes on purpose with UnmarshalCSVSample's Timestamp field. Since
	// *UnmarshalCSVSample implements UnmarshalCSV(), that method call should
	// take precedence.
	Timestamp *UnmarshalCSVSample
}

type EmbedPtrSample struct {
	Qux string `csv:"first"`
	*Sample
	Ignore string  `csv:"-"`
	Grault float64 `csv:"garply"`
	Quux   string  `csv:"last"`
}

type SkipFieldSample struct {
	EmbedSample
	MoreIgnore string `csv:"-"`
	Corge      string `csv:"abc"`
}

// Testtype for unmarshal/marshal functions on renamed basic types
type RenamedFloat64Unmarshaler float64
type RenamedFloat64Default float64

type RenamedSample struct {
	RenamedFloatUnmarshaler RenamedFloat64Unmarshaler `csv:"foo"`
	RenamedFloatDefault     RenamedFloat64Default     `csv:"bar"`
}

type MultiTagSample struct {
	Foo string `csv:"Baz,foo"`
	Bar int    `csv:"BAR"`
}

type TagSeparatorSample struct {
	Foo string `csv:"Baz|foo"`
	Bar int    `csv:"BAR"`
}

type CustomTagSample struct {
	Foo string `custom:"foo"`
	Bar string `csv:"BAR"`
}

type DateTime struct {
	Foo time.Time `csv:"Foo"`
}

type Level0Struct struct {
	Level0Field level1Struct
}

type level1Struct struct {
	Level1Field level2Struct
}

type level2Struct struct {
	InnerStruct
}

type InnerStruct struct {
	BoolIgnoreField0 bool   `csv:"-"`
	BoolField1       bool   `csv:"boolField1"`
	StringField2     string `csv:"stringField2"`
}

var _ TypeUnmarshalCSVWithFields = (*UnmarshalCSVWithFieldsSample)(nil)

type UnmarshalCSVWithFieldsSample struct {
	Foo  string  `csv:"foo"`
	Bar  int     `csv:"bar"`
	Baz  string  `csv:"baz"`
	Frop float64 `csv:"frop"`
}
