package gocsv

import (
	"encoding/csv"
	"strings"
	"testing"
)

type CSVDate struct {
	Date string
}

func (self *CSVDate) UnmarshalCSV(text string) error {
	if self == nil {
		self = &CSVDate{}
	}
	self.Date = text
	return nil
}

func Test_CSV_Base(t *testing.T) {
	t.Parallel()

	type row struct {
		ID   string   `csv:"id"`
		Date *CSVDate `csv:"date"`
	}

	exampleCSV := `id,date
1,foo
2,bar
`

	var rows []row
	r := csv.NewReader(strings.NewReader(exampleCSV))
	err := UnmarshalCSV(r, &rows)
	if err != nil {
		t.Fatal(err.Error())
	}

	if rows[0].Date.Date != "foo" {
		t.Fatalf("Expected %q, but got %q", "foo", string(rows[0].Date.Date))
	}
}

////////////////////////////////////////////////////////////

type FieldWithCustomMarshaller struct {
	value string
}

func (f *FieldWithCustomMarshaller) UnmarshalCSV(csv string) (err error) {
	f.value = csv
	return err
}

type FieldWithCustomMarshallerPointed struct {
	otherValue string
}

func (f *FieldWithCustomMarshallerPointed) UnmarshalCSV(csv string) (err error) {
	f.otherValue = csv
	return err
}

type TestStruct struct {
	OkValue                          string
	FieldWithCustomMarshaller        FieldWithCustomMarshaller
	FieldWithCustomMarshallerPointed *FieldWithCustomMarshallerPointed
}

func TestPanic(t *testing.T) {
	line := "make,backups,test it"
	r := strings.NewReader(line)
	var DataValues []TestStruct
	err := UnmarshalWithoutHeaders(r, &DataValues)
	if err != nil {
		t.Fatal(err)
	}
}
