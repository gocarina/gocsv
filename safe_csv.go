package gocsv

//Wraps around SafeCSVWriter and makes it thread safe.
import (
	"encoding/csv"
	"sync"
)

//CSVWriter interface for anything implementing csv writing api
type CSVWriter interface {
	Write(row []string) error
	Flush()
	Error() error
}

//SafeCSVWriter mutex protected thread safe csv writer
type SafeCSVWriter struct {
	*csv.Writer
	m sync.Mutex
}

//NewSafeCSVWriter create a new SafeCSVWriter
func NewSafeCSVWriter(original *csv.Writer) *SafeCSVWriter {
	return &SafeCSVWriter{
		Writer: original,
	}
}

//Write the csv writer in a threadsafe way
func (w *SafeCSVWriter) Write(row []string) error {
	w.m.Lock()
	defer w.m.Unlock()
	return w.Writer.Write(row)
}

//Flush flush the csv writer in a threadsafe way
func (w *SafeCSVWriter) Flush() {
	w.m.Lock()
	w.Writer.Flush()
	w.m.Unlock()
}
