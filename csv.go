// Copyright 2014 Jonathan Picques. All rights reserved.
// Use of this source code is governed by a MIT license
// The license can be found in the LICENSE file.

// The GoCSV package aims to provide easy CSV serialization and deserialization to the golang programming language

package gocsv

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"strings"
)

var FailIfUnmatchedStructTags = false

// --------------------------------------------------------------------------
// CSVWriter used to format CSV

var selfCSVWriter = DefaultCSVWriter

// DefaultCSVWriter is the default CSV writer used to format CSV (cf. csv.NewWriter)
func DefaultCSVWriter(out io.Writer) *csv.Writer {
	return csv.NewWriter(out)
}

// SetCSVWriter sets the CSV writer used to format CSV.
func SetCSVWriter(csvWriter func(io.Writer) *csv.Writer) {
	selfCSVWriter = csvWriter
}

func getCSVWriter(out io.Writer) *csv.Writer {
	return selfCSVWriter(out)
}

// --------------------------------------------------------------------------
// CSVReader used to parse CSV

var selfCSVReader = DefaultCSVReader

// DefaultCSVReader is the default CSV reader used to parse CSV (cf. csv.NewReader)
func DefaultCSVReader(in io.Reader) *csv.Reader {
	return csv.NewReader(in)
}

// LazyCSVReader returns a lazy CSV reader, with LazyQuotes and TrimLeadingSpace.
func LazyCSVReader(in io.Reader) *csv.Reader {
	csvReader := csv.NewReader(in)
	csvReader.LazyQuotes = true
	csvReader.TrimLeadingSpace = true
	return csvReader
}

// SetCSVReader sets the CSV reader used to parse CSV.
func SetCSVReader(csvReader func(io.Reader) *csv.Reader) {
	selfCSVReader = csvReader
}

func getCSVReader(in io.Reader) *csv.Reader {
	return selfCSVReader(in)
}

// --------------------------------------------------------------------------
// Marshal functions

// MarshalFile saves the interface as CSV in the file.
func MarshalFile(in interface{}, file *os.File) (err error) {
	return Marshal(in, file)
}

// MarshalString returns the CSV string from the interface.
func MarshalString(in interface{}) (out string, err error) {
	bufferString := bytes.NewBufferString(out)
	if err := Marshal(in, bufferString); err != nil {
		return "", err
	}
	return bufferString.String(), nil
}

// MarshalBytes returns the CSV bytes from the interface.
func MarshalBytes(in interface{}) (out []byte, err error) {
	bufferString := bytes.NewBuffer(out)
	if err := Marshal(in, bufferString); err != nil {
		return nil, err
	}
	return bufferString.Bytes(), nil
}

// Marshal returns the CSV in writer from the interface.
func Marshal(in interface{}, out io.Writer) (err error) {
	writer := getCSVWriter(out)
	return writeTo(writer, in)
}

func MarshalCSV(in interface{}, out *csv.Writer) (err error) {
	return writeTo(out, in)
}

// --------------------------------------------------------------------------
// Unmarshal functions

// UnmarshalFile parses the CSV from the file in the interface.
func UnmarshalFile(in *os.File, out interface{}) (err error) {
	return Unmarshal(in, out)
}

// UnmarshalString parses the CSV from the string in the interface.
func UnmarshalString(in string, out interface{}) (err error) {
	return Unmarshal(strings.NewReader(in), out)
}

// UnmarshalBytes parses the CSV from the bytes in the interface.
func UnmarshalBytes(in []byte, out interface{}) (err error) {
	return Unmarshal(bytes.NewReader(in), out)
}

// Unmarshal parses the CSV from the reader in the interface.
func Unmarshal(in io.Reader, out interface{}) (err error) {
	return readTo(newDecoder(in), out)
}

func UnmarshalCSV(in *csv.Reader, out interface{}) error {
	return readTo(csvDecoder{in}, out)
}
