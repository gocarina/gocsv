package gocsv

//Wraps around SafeCSVWriter and makes it thread safe.
import (
	"sync"
)

type CSVWriter interface {
	Write(row []string) error
	Flush()
	Error() error
}

type SafeCSVWriter struct {
	CSVWriter
	m sync.Mutex
}

func NewSafeCSVWriter(original CSVWriter) *SafeCSVWriter {
	return &SafeCSVWriter{
		CSVWriter: original,
	}
}

//Override write
func (w *SafeCSVWriter) Write(row []string) error {
	w.m.Lock()
	defer w.m.Unlock()
	return w.CSVWriter.Write(row)
}

//Override flush
func (w *SafeCSVWriter) Flush() {
	w.m.Lock()
	w.CSVWriter.Flush()
	w.m.Unlock()
}
