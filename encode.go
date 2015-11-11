package gocsv

import (
	"fmt"
	"io"
	"reflect"
)

type encoder struct {
	out io.Writer
}

func newEncoder(out io.Writer) *encoder {
	return &encoder{out}
}

func (encode *encoder) writeTo(in interface{}) error {
	inValue, inType := getConcreteReflectValueAndType(in) // Get the concrete type (not pointer) (Slice<?> or Array<?>)
	if err := encode.ensureInType(inType); err != nil {
		return err
	}
	inInnerWasPointer, inInnerType := getConcreteContainerInnerType(inType) // Get the concrete inner type (not pointer) (Container<"?">)
	if err := encode.ensureInInnerType(inInnerType); err != nil {
		return err
	}
	csvWriter := getCSVWriter(encode.out)           // Get the CSV writer
	inInnerStructInfo := getStructInfo(inInnerType) // Get the inner struct info to get CSV annotations
	csvHeadersLabels := make([]string, len(inInnerStructInfo.Fields))
	for i, fieldInfo := range inInnerStructInfo.Fields { // Used to write the header (first line) in CSV
		csvHeadersLabels[i] = fieldInfo.Key
	}
	csvWriter.Write(csvHeadersLabels)
	inLen := inValue.Len()
	for i := 0; i < inLen; i++ { // Iterate over container rows
		for j, fieldInfo := range inInnerStructInfo.Fields {
			csvHeadersLabels[j] = ""
			inInnerFieldValue, err := encode.getInnerField(inValue.Index(i), inInnerWasPointer, fieldInfo.Num) // Get the correct field header <-> position
			if err != nil {
				return err
			}
			csvHeadersLabels[j] = inInnerFieldValue
		}
		csvWriter.Write(csvHeadersLabels)
	}
	csvWriter.Flush()
	return csvWriter.Error()
}

// Check if the inType is an array or a slice
func (encode *encoder) ensureInType(outType reflect.Type) error {
	switch outType.Kind() {
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		return nil
	}
	return fmt.Errorf("cannot use " + outType.String() + ", only slice or array supported")
}

// Check if the inInnerType is of type struct
func (encode *encoder) ensureInInnerType(outInnerType reflect.Type) error {
	switch outInnerType.Kind() {
	case reflect.Struct:
		return nil
	}
	return fmt.Errorf("cannot use " + outInnerType.String() + ", only struct supported")
}

func (encode *encoder) getInnerField(outInner reflect.Value, outInnerWasPointer bool, fieldPosition int) (string, error) {
	if outInnerWasPointer {
		return getFieldAsString(outInner.Elem().Field(fieldPosition))
	}
	return getFieldAsString(outInner.Field(fieldPosition))
}
