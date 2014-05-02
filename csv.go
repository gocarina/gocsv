// Copyright 2014 Jonathan Picques. All rights reserved.
// Use of this source code is governed by a MIT license
// The license can be found in the LICENSE file.

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

var selfCSVWriter func(io.Writer) *csv.Writer = DefaultCSVWriter

// Default CSV writer (see csv.NewWriter)
func DefaultCSVWriter(out io.Writer) *csv.Writer {
	return csv.NewWriter(out)
}

// Set the CSV writer used to unmarshal.
func SetCSVWriter(csvWriter func(io.Writer) *csv.Writer) {
	selfCSVWriter = csvWriter
}

func getCSVWriter(out io.Writer) *csv.Writer {
	return selfCSVWriter(out)
}

// --------------------------------------------------------------------------
// CSVReader used to parse CSV

var selfCSVReader func(io.Reader) *csv.Reader = DefaultCSVReader

// Default CSV reader (see csv.NewReader)
func DefaultCSVReader(in io.Reader) *csv.Reader {
	return csv.NewReader(in)
}

// Get a lazy CSV reader, with LazyQuotes and TrimLeadingSpace.
func LazyCSVReader(in io.Reader) *csv.Reader {
	csvReader := csv.NewReader(in)
	csvReader.LazyQuotes = true
	csvReader.TrimLeadingSpace = true
	return csvReader
}

// Set the CSV reader used to marshal.
func SetCSVReader(csvReader func(io.Reader) *csv.Reader) {
	selfCSVReader = csvReader
}

func getCSVReader(in io.Reader) *csv.Reader {
	return selfCSVReader(in)
}

// --------------------------------------------------------------------------
// Marshal functions

// Save the in interface as CSV in file.
func MarshalFile(in interface{}, file *os.File) (err error) {
	return Marshal(in, file)
}

// Returns the CSV string from in interface.
func MarshalString(in interface{}) (out string, err error) {
	bufferString := bytes.NewBufferString(out)
	if err := Marshal(in, bufferString); err != nil {
		return "", err
	}
	return bufferString.String(), nil
}

// Returns the CSV bytes from in interface.
func MarshalBytes(in interface{}) (out []byte, err error) {
	bufferString := bytes.NewBuffer(out)
	if err := Marshal(in, bufferString); err != nil {
		return nil, err
	}
	return bufferString.Bytes(), nil
}

// Returns the CSV in writer from in interface
func Marshal(in interface{}, out io.Writer) (err error) {
	return newEncoder(out).writeTo(in)
}

// --------------------------------------------------------------------------
// Unmarshal functions

// Unmarshal the file in out interface.
func UnmarshalFile(in *os.File, out interface{}) (err error) {
	return Unmarshal(in, out)
}

// Unmarshal the string in out interface.
func UnmarshalString(in string, out interface{}) (err error) {
	return Unmarshal(strings.NewReader(in), out)
}

// Unmarshal the bytes in out interface
func UnmarshalBytes(in []byte, out interface{}) (err error) {
	return Unmarshal(bytes.NewReader(in), out)
}

// Unmarshal the data from reader in out interface
func Unmarshal(in io.Reader, out interface{}) (err error) {
	return newDecoder(in).readTo(out)
}

