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

func (self *encoder) writeTo(in reflect.Value) error {
	if err := self.ensureInKind(&in); err != nil {
		return err
	}
	inType := in.Type()
	inInnerType := inType.Elem()
	if err := self.ensureInInnerKind(inInnerType); err != nil { // Check if internal data is a struct
		return err
	}
	inInnerTypeStructInfo := getStructInfo(inInnerType)
	csvColumnsNames := make([]string, len(inInnerTypeStructInfo.Fields))
	for i, inInnerField := range inInnerTypeStructInfo.Fields {
		csvColumnsNames[i] = inInnerField.Key
	}
	csvWriter := getCSVWriter(self.out)
	csvWriter.Write(csvColumnsNames)
	outLen := in.Len()
	for i := 0; i < outLen; i++ {
		for j, inInnerField := range inInnerTypeStructInfo.Fields {
			inInnerFieldValue, err := getFieldAsString(in.Index(i).Field(inInnerField.Num))
			if err != nil {
				return err
			}
			csvColumnsNames[j] = inInnerFieldValue
		}
		csvWriter.Write(csvColumnsNames)
	}
	csvWriter.Flush()
	return nil
}

func (self *encoder) ensureInKind(in *reflect.Value) error {
	switch in.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		return nil
	}
	return fmt.Errorf("Unsupported type " + in.Type().String() + ", only Slice or Array supported")
}

func (self *encoder) ensureInInnerKind(inInnerType reflect.Type) error {
	switch inInnerType.Kind() {
	case reflect.Struct:
		return nil
	}
	return fmt.Errorf("Unsupported type " + inInnerType.String() + ", only struct supported")
}
