package gocsv_test

import (
	"encoding/csv"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/gocarina/gocsv"
)

// heightDecoder normalizes only Height values after CSV parsing. Keeping the
// header unchanged also keeps GoCSV's ordinary field matching in effect.
type heightDecoder struct {
	reader *csv.Reader
}

func (d heightDecoder) GetCSVRows() ([][]string, error) {
	rows, err := d.reader.ReadAll()
	if err != nil || len(rows) == 0 {
		return rows, err
	}
	for column, header := range rows[0] {
		if header != "Height" {
			continue
		}
		for _, row := range rows[1:] {
			if row[column] == "N/A" {
				row[column] = "NaN"
			}
		}
	}
	return rows, nil
}

type measuredPerson struct {
	Name   string  `csv:"Name"`
	Age    int     `csv:"Age"`
	Height float64 `csv:"Height"`
}

func ExampleUnmarshalDecoder() {
	input := "Name,Age,Height\nJohn,10,1.80\n\"Jack, Jr.\",20,N/A\nZoë,30,1.70\n"
	decoder := heightDecoder{reader: csv.NewReader(strings.NewReader(input))}
	var people []measuredPerson
	if err := gocsv.UnmarshalDecoder(decoder, &people); err != nil {
		fmt.Println(err)
		return
	}
	for _, person := range people {
		fmt.Printf("%s %d %v\n", person.Name, person.Age, person.Height)
	}
	// Output:
	// John 10 1.8
	// Jack, Jr. 20 NaN
	// Zoë 30 1.7
}

func TestSanitizingDecoderControls(t *testing.T) {
	// Only the configured column is mapped, even when its position changes.
	var people []measuredPerson
	decoder := heightDecoder{reader: csv.NewReader(strings.NewReader("Height,Name,Age\nN/A,N/A,20\n1.80,Zoë,30\n"))}
	if err := gocsv.UnmarshalDecoder(decoder, &people); err != nil {
		t.Fatal(err)
	}
	if len(people) != 2 || people[0].Name != "N/A" || !math.IsNaN(people[0].Height) || people[1].Height != 1.8 {
		t.Fatalf("unexpected mapped records: %#v", people)
	}

	for _, input := range []string{
		"Name,Age,Height\nJohn,10\n",         // Malformed CSV is not repaired.
		"Name,Age,Height\nJohn,10,unknown\n", // Other numeric errors remain errors.
		"Name,Age,Height\nJohn,N/A,1.80\n",   // Other columns are not mapped.
	} {
		t.Run(input, func(t *testing.T) {
			var records []measuredPerson
			decoder := heightDecoder{reader: csv.NewReader(strings.NewReader(input))}
			if err := gocsv.UnmarshalDecoder(decoder, &records); err == nil {
				t.Fatal("expected an error")
			}
		})
	}

	// Without the application-specific decoder, the sentinel is invalid for a
	// float, as expected: GoCSV cannot infer what an arbitrary sentinel means.
	if err := gocsv.UnmarshalString("Name,Age,Height\nJack,20,N/A\n", &people); err == nil {
		t.Fatal("expected normal float conversion to reject N/A")
	}
}
