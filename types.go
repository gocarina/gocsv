package gocsv

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// --------------------------------------------------------------------------
// Reflection helpers

type structInfo struct {
	Fields []fieldInfo
	Zero   reflect.Value
}

type fieldInfo struct {
	Key       string
	Num       int
	OmitEmpty bool
}

var structMap = make(map[reflect.Type]*structInfo)
var structMapMutex sync.RWMutex

func getInterfaceType(out interface{}) reflect.Value {
	outReflect := reflect.ValueOf(out)
	outReflectKind := outReflect.Kind()

	for outReflectKind == reflect.Ptr {
		outReflect = outReflect.Elem()
		outReflectKind = outReflect.Kind()
	}
	return outReflect
}

func getStructInfo(rType reflect.Type) *structInfo {
	structMapMutex.RLock()
	stInfo, ok := structMap[rType]
	structMapMutex.RUnlock()
	if ok {
		return stInfo
	}
	fieldsCount := rType.NumField()
	fieldsList := make([]fieldInfo, 0, fieldsCount)
	for i := 0; i < fieldsCount; i++ {
		field := rType.Field(i)
		if field.PkgPath != "" {
			continue
		}
		fieldInfo := fieldInfo{Num: i}
		fieldTag := field.Tag.Get("csv")
		fieldTags := strings.Split(fieldTag, ",")
		for _, fieldTagEntry := range fieldTags {
			if fieldTagEntry == "omitempty" {
				fieldInfo.OmitEmpty = true
			} else {
				fieldTag = fieldTagEntry
			}
		}
		if fieldTag == "-" {
			continue
		} else if fieldTag != "" {
			fieldInfo.Key = fieldTag
		} else {
			fieldInfo.Key = field.Name
		}
		fieldsList = append(fieldsList, fieldInfo)
	}
	stInfo = &structInfo{fieldsList, reflect.New(rType).Elem()}
	return stInfo
}

// --------------------------------------------------------------------------
// Conversion interfaces

type TypeMarshaller interface {
	MarshalCSV() (string, error)
}

type TypeUnmarshaller interface {
	UnmarshalCSV(string) error
}

// --------------------------------------------------------------------------
// Conversion helpers

func toString(in interface{}) (string, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		return inValue.String(), nil
	case reflect.Bool:
		b := inValue.Bool()
		if b {
			return "true", nil
		} else {
			return "false", nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return string(inValue.Int()), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return string(inValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(inValue.Float(), byte('f'), 64, 64), nil
	}
	return "", fmt.Errorf("No known conversion from " + inValue.Type().String() + " to string")
}

func toBool(in interface{}) (bool, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := inValue.String()
		if s == "true" || s == "yes" || s == "1" {
			return true, nil
		} else if s == "false" || s == "no" || s == "0" {
			return false, nil
		}
	case reflect.Bool:
		return inValue.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := inValue.Int()
		if i != 0 {
			return true, nil
		} else {
			return false, nil
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := inValue.Uint()
		if i != 0 {
			return true, nil
		} else {
			return false, nil
		}
	case reflect.Float32, reflect.Float64:
		f := inValue.Float()
		if f != 0 {
			return true, nil
		} else {
			return false, nil
		}
	}
	return false, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to bool")
}

func toInt(in interface{}) (int64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		return strconv.ParseInt(inValue.String(), 0, 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		} else {
			return 0, nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return inValue.Int(), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(inValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(inValue.Float()), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to int")
}

func toUint(in interface{}) (uint64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		return strconv.ParseUint(inValue.String(), 0, 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		} else {
			return 0, nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(inValue.Int()), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return inValue.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return uint64(inValue.Float()), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to uint")
}

func toFloat(in interface{}) (float64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		return strconv.ParseFloat(inValue.String(), 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		} else {
			return 0, nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(inValue.Int()), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(inValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return inValue.Float(), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to float")
}

func setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		s, err := toString(value)
		if err != nil {
			return err
		}
		field.SetString(s)
	case reflect.Bool:
		b, err := toBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := toInt(value)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ui, err := toUint(value)
		if err != nil {
			return err
		}
		field.SetUint(ui)
	case reflect.Float32, reflect.Float64:
		f, err := toFloat(value)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	default:
		saveField := field
		if field.Kind() == reflect.Ptr && field.IsNil() {
			field = reflect.New(field.Type().Elem())
		}
		in, ok := field.Interface().(TypeUnmarshaller)
		if !ok {
			return fmt.Errorf("No known conversion from " + reflect.TypeOf(value).String() + " to " + field.Type().String() + ", " + field.Type().String() + " does not implements TypeUnmarshaller")
		}
		if err := in.UnmarshalCSV(value); err != nil {
			return err
		}
		saveField.Set(reflect.ValueOf(in))
	}
	return nil
}

func getFieldAsString(field reflect.Value) (str string, err error) {
	switch field.Kind() {
	case reflect.String:
		return field.String(), nil
	case reflect.Bool:
		str, err = toString(field.Bool())
		if err != nil {
			return str, err
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str, err = toString(field.Int())
		if err != nil {
			return str, err
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str, err = toString(field.Uint())
		if err != nil {
			return str, err
		}
	case reflect.Float32, reflect.Float64:
		str, err = toString(field.Float())
		if err != nil {
			return str, err
		}
	default:
		if field.Kind() == reflect.Ptr && field.IsNil() {
			field = reflect.New(field.Type().Elem())
		}
		out, ok := field.Interface().(TypeMarshaller)
		if !ok {
			return str, fmt.Errorf("No known conversion from " + field.Type().String() + " to string, " + field.Type().String() + " does not implements TypeMarshaller")
		}
		str, err = out.MarshalCSV()
		if err != nil {
			return str, err
		}
	}
	return str, nil
}
