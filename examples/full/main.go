package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type NotUsed struct {
	Name string
}

type Client struct { // Our example struct, you can use "-" to ignore a field
	ID            string   `csv:"client_id"`
	Name          string   `csv:"client_name"`
	Age           string   `csv:"client_age"`
	NotUsedString string   `csv:"-"`
	NotUsedStruct NotUsed  `csv:"-"`
	Address1      Address  `csv:"addr1"`
	Address2      Address  //`csv:"addr2"` will use Address2 in header
	Employed      DateTime `csv:"employed"`
}

type Address struct {
	Street string `csv:"street"`
	City   string `csv:"city"`
}

var _ gocsv.TypeMarshaller = new(DateTime)
var _ gocsv.TypeUnmarshaller = new(DateTime)

type DateTime struct {
	time.Time
}

// Convert the internal date as CSV string
func (date *DateTime) MarshalCSV() (string, error) {
	return date.String(), nil
}

// You could also use the standard Stringer interface
func (date DateTime) String() string {
	return date.Time.Format("20060201")
}

// Convert the CSV string as internal date
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("20060201", csv)
	return err
}

func main() {
	// set the pipe as the delimiter for writing
	gocsv.TagSeparator = "|"
	// set the pipe as the delimiter for reading
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '|'
		return r
	})

	// Create an empty clients file
	clientsFile, err := os.OpenFile("clients.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer clientsFile.Close()

	// Create clients
	clients := []*Client{
		{ID: "12", Name: "John", Age: "21",
			Address1: Address{"Street 1", "City1"},
			Address2: Address{"Street 2", "City2"},
			Employed: DateTime{time.Date(2022, 11, 04, 12, 0, 0, 0, time.UTC)},
		},
		{ID: "13", Name: "Fred",
			Address1: Address{`Main "Street" 1`, "City1"}, // show quotes in value
			Address2: Address{"Main Street 2", "City2"},
			Employed: DateTime{time.Date(2022, 11, 04, 13, 0, 0, 0, time.UTC)},
		},
		{ID: "14", Name: "James", Age: "32",
			Address1: Address{"Center Street 1", "City1"},
			Address2: Address{"Center Street 2", "City2"},
			Employed: DateTime{time.Date(2022, 11, 04, 14, 0, 0, 0, time.UTC)},
		},
		{ID: "15", Name: "Danny",
			Address1: Address{"State Street 1", "City1"},
			Address2: Address{"State Street 2", "City2"},
			Employed: DateTime{time.Date(2022, 11, 04, 15, 0, 0, 0, time.UTC)},
		},
	}
	// Save clients to csv file
	if err = gocsv.MarshalFile(&clients, clientsFile); err != nil {
		panic(err)
	}

	// Reset the file reader
	if _, err := clientsFile.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}
	// Read file and print to console
	b, err := ioutil.ReadAll(clientsFile)
	if err != nil {
		panic(err)
	}
	fmt.Println("clients.csv:")
	fmt.Println(string(b))

	// Reset the file reader
	if _, err := clientsFile.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}
	// Load clients from file
	var newClients []*Client
	if err := gocsv.UnmarshalFile(clientsFile, &newClients); err != nil {
		panic(err)
	}
	fmt.Println("clients:")
	for _, c := range newClients {
		fmt.Printf("%s:%s Adress1:%s, %s Address2:%s, %s Employed: %s\n",
			c.ID, c.Name,
			c.Address1.Street, c.Address1.City,
			c.Address2.Street, c.Address2.City,
			c.Employed,
		)
	}
}
