package gocsv

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"strings"
)

// --------------------------------------------------------------------------
// CSVReader used to format CSV

type CSVWriter func(io.Writer) *csv.Writer

var selfCSVWriter CSVWriter = DefaultCSVWriter

func DefaultCSVWriter(out io.Writer) *csv.Writer {
	return csv.NewWriter(out)
}

func SetCSVWriter(csvWriter CSVWriter) {
	selfCSVWriter = csvWriter
}

func getCSVWriter(out io.Writer) *csv.Writer {
	return selfCSVWriter(out)
}

// --------------------------------------------------------------------------
// CSVReader used to parse CSV

type CSVReader func(io.Reader) *csv.Reader

var selfCSVReader CSVReader = DefaultCSVReader

func DefaultCSVReader(in io.Reader) *csv.Reader {
	return csv.NewReader(in)
}

func LazyCSVReader(in io.Reader) *csv.Reader {
	csvReader := csv.NewReader(in)
	csvReader.LazyQuotes = true
	csvReader.TrimLeadingSpace = true
	return csvReader
}

func SetCSVReader(csvReader func(io.Reader) *csv.Reader) {
	selfCSVReader = csvReader
}

func getCSVReader(in io.Reader) *csv.Reader {
	return selfCSVReader(in)
}

// --------------------------------------------------------------------------
// Marshal functions

func MarshalFile(in interface{}, file *os.File) (err error) {
	return Marshal(in, file)
}

func MarshalString(in interface{}) (out string, err error) {
	bufferString := bytes.NewBufferString(out)
	if err := Marshal(in, bufferString); err != nil {
		return "", err
	}
	return bufferString.String(), nil
}

func MarshalBytes(in interface{}) (out []byte, err error) {
	bufferString := bytes.NewBuffer(out)
	if err := Marshal(in, bufferString); err != nil {
		return nil, err
	}
	return bufferString.Bytes(), nil
}

func Marshal(in interface{}, out io.Writer) (err error) {
	return newEncoder(out).writeTo(getInterfaceType(in))
}

// --------------------------------------------------------------------------
// Unmarshal functions

func UnmarshalFile(in *os.File, out interface{}) (err error) {
	return Unmarshal(in, out)
}

func UnmarshalString(in string, out interface{}) (err error) {
	return Unmarshal(strings.NewReader(in), out)
}

func UnmarshalBytes(in []byte, out interface{}) (err error) {
	return Unmarshal(bytes.NewReader(in), out)
}

func Unmarshal(in io.Reader, out interface{}) (err error) {
	return newDecoder(in).readTo(getInterfaceType(out))
}

