package gocsv

import (
	"fmt"
	"io"
	"reflect"
)

type decoder struct {
	in io.Reader
}

func newDecoder(in io.Reader) *decoder {
	return &decoder{in}
}

func (self *decoder) readTo(out reflect.Value) error {
	if err := self.ensureOutKind(&out); err != nil { // Check if interface is of type Slice or Array
		return err
	}
	outType := out.Type()
	csvContent, err := self.getCSVContent() // Get the CSV content
	if err != nil {
		return err
	}
	if err := self.ensureOutCapacity(&out, len(csvContent)); err != nil { // Check capacity and grows it if possible
		return err
	}
	outInnerType := outType.Elem()
	if err := self.ensureOutInnerKind(outInnerType); err != nil { // Check if internal data is a struct
		return err
	}
	outInnerStructInfo := getStructInfo(outInnerType) // Get struct info to get the columns tags
	csvColumnsNamesFieldsIndex := make(map[int]int)   // Used to store column names and position in struct
	for i, csvRow := range csvContent {

		if i == 0 {
			for j, csvColumn := range csvRow {
				if num := self.ensureCSVFieldExists(csvColumn, outInnerStructInfo); num != -1 {
					csvColumnsNamesFieldsIndex[j] = num
				}
			}
		} else {
			outInner := reflect.New(outInnerType).Elem()
			for j, csvColumn := range csvRow {
				if pos, ok := csvColumnsNamesFieldsIndex[j]; ok { // Position found accordingly to column name position
					if err := setField(outInner.Field(pos), csvColumn); err != nil { // Set field of struct
						return err
					}
				}
			}
			out.Index(i - 1).Set(outInner) // Position if offset by one (0 is column names)
		}
	}
	return nil
}

func (self *decoder) ensureOutKind(out *reflect.Value) error {
	switch out.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		return nil
	}
	return fmt.Errorf("Unsupported type " + out.Type().String() + ", only slice or array supported")
}

func (self *decoder) ensureOutInnerKind(outInnerType reflect.Type) error {
	switch outInnerType.Kind() {
	case reflect.Struct:
		return nil
	}
	return fmt.Errorf("Unsupported type " + outInnerType.String() + ", only struct supported")
}

func (self *decoder) ensureOutCapacity(out *reflect.Value, csvLen int) error {

	switch out.Kind() {
	case reflect.Array:
		if out.Cap() < csvLen-1 { // Array is not big enough to hold the CSV content
			return fmt.Errorf("Array capacity problem")
		}
	case reflect.Slice:
		if !out.CanAddr() && out.Len() < csvLen-1 { // Slice is not big enough tho hold the CSV content and is not addressable
			return fmt.Errorf("Slice not addressable capacity problem (did you forget &?)")
		} else if out.CanAddr() && out.Len() < csvLen-1 {
			out.Set(reflect.MakeSlice(out.Type(), csvLen-1, csvLen-1)) // Slice is not big enough, so grows it
		}
	}
	return nil
}

func (self *decoder) ensureCSVFieldExists(key string, structInfo *structInfo) int {
	for i, field := range structInfo.Fields {
		if field.Key == key {
			return i
		}
	}
	return -1
}

func (self *decoder) getCSVContent() ([][]string, error) {
	return getCSVReader(self.in).ReadAll()
}
