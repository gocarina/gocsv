package gocsv

import (
	"reflect"
	"strings"
	"sync"
)

// --------------------------------------------------------------------------
// Reflection helpers

type structInfo struct {
	Fields []fieldInfo
}

type fieldInfo struct {
	Key string
	Num int
}

var structMap = make(map[reflect.Type]*structInfo)
var structMapMutex sync.RWMutex

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
			if fieldTagEntry != "omitempty" {
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
	stInfo = &structInfo{fieldsList}
	return stInfo
}

func getConcreteContainerInnerType(in reflect.Type) (inInnerWasPointer bool, inInnerType reflect.Type) {
	inInnerType = in.Elem()
	inInnerWasPointer = false
	if inInnerType.Kind() == reflect.Ptr {
		inInnerWasPointer = true
		inInnerType = inInnerType.Elem()
	}
	return inInnerWasPointer, inInnerType
}

func getConcreteReflectValueAndType(in interface{}) (reflect.Value, reflect.Type) {
	value := reflect.ValueOf(in)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value, value.Type()
}
