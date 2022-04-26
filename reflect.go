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
}

// fieldInfo is a struct field that should be mapped to a CSV column, or vice-versa
// Each IndexChain element before the last is the index of an the embedded struct field
// that defines Key as a tag
type fieldInfo struct {
	keys         []string
	omitEmpty    bool
	IndexChain   []int
	defaultValue string
}

func (f fieldInfo) getFirstKey() string {
	return f.keys[0]
}

func (f fieldInfo) matchesKey(key string) bool {
	for _, k := range f.keys {
		if key == k || strings.TrimSpace(key) == k {
			return true
		}
	}
	return false
}

var structInfoCache sync.Map
var structMap = make(map[reflect.Type]*structInfo)
var structMapMutex sync.RWMutex

func getStructInfo(rType reflect.Type) *structInfo {
	stInfo, ok := structInfoCache.Load(rType)
	if ok {
		return stInfo.(*structInfo)
	}

	fieldsList := getFieldInfos(rType, []int{})
	stInfo = &structInfo{fieldsList}
	structInfoCache.Store(rType, stInfo)

	return stInfo.(*structInfo)
}

func getFieldInfos(rType reflect.Type, parentIndexChain []int) []fieldInfo {
	fieldsCount := rType.NumField()
	fieldsList := make([]fieldInfo, 0, fieldsCount)
	for i := 0; i < fieldsCount; i++ {
		field := rType.Field(i)
		if field.PkgPath != "" {
			continue
		}

		var cpy = make([]int, len(parentIndexChain))
		copy(cpy, parentIndexChain)
		indexChain := append(cpy, i)
		// if the field is a pointer to a struct, follow the pointer then create fieldinfo for each field
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			// unless it implements marshalText or marshalCSV. Structs that implement this
			// should result in one value and not have their fields exposed
			if !(canMarshal(field.Type.Elem())) {
				fieldsList = append(fieldsList, getFieldInfos(field.Type.Elem(), indexChain)...)
			}
		}
		// if the field is a struct, create a fieldInfo for each of its fields
		if field.Type.Kind() == reflect.Struct {
			// unless it implements marshalText or marshalCSV. Structs that implement this
			// should result in one value and not have their fields exposed
			if !(canMarshal(field.Type)) {
				fieldsList = append(fieldsList, getFieldInfos(field.Type, indexChain)...)
			}
		}

		// if the field is an embedded struct, ignore the csv tag
		if field.Anonymous {
			continue
		}

		currFieldInfo := fieldInfo{IndexChain: indexChain}
		fieldTag := field.Tag.Get(TagName)
		fieldTags := strings.Split(fieldTag, TagSeparator)
		filteredTags := []string{}
		for _, fieldTagEntry := range fieldTags {
			trimmedFieldTagEntry := strings.TrimSpace(fieldTagEntry) // handles cases like `csv:"foo, omitempty, default=test"`
			if trimmedFieldTagEntry == "omitempty" {
				currFieldInfo.omitEmpty = true
			} else if strings.HasPrefix(trimmedFieldTagEntry, "default=") {
				currFieldInfo.defaultValue = strings.TrimPrefix(trimmedFieldTagEntry, "default=")
			} else {
				filteredTags = append(filteredTags, normalizeName(trimmedFieldTagEntry))
			}
		}

		if len(filteredTags) == 1 && filteredTags[0] == "-" {
			continue
		} else if len(filteredTags) > 0 && filteredTags[0] != "" {
			currFieldInfo.keys = filteredTags
		} else {
			currFieldInfo.keys = []string{normalizeName(field.Name)}
		}

		if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array {
			var arrayLength = -1
			if arrayTag, ok := field.Tag.Lookup(TagName + "[]"); ok {
				arrayLength, _ = strconv.Atoi(arrayTag)
			}

			// When the field is a slice/array of structs, create a fieldInfo for each index and each field
			if field.Type.Elem().Kind() == reflect.Struct {
				fieldInfos := getFieldInfos(field.Type.Elem(), []int{})

				for idx := 0; idx < arrayLength; idx++ {
					// copy index chain and append array index
					var cpy2 = make([]int, len(indexChain))
					copy(cpy2, indexChain)
					arrayIndexChain := append(cpy2, idx)
					for _, childFieldInfo := range fieldInfos {
						// copy array index chain and append array index
						var cpy3 = make([]int, len(arrayIndexChain))
						copy(cpy3, arrayIndexChain)

						arrayFieldInfo := fieldInfo{
							IndexChain:   append(cpy3, childFieldInfo.IndexChain...),
							omitEmpty:    childFieldInfo.omitEmpty,
							defaultValue: childFieldInfo.defaultValue,
						}

						// create cartesian product of keys
						// eg: array field keys x struct field keys
						for _, akey := range currFieldInfo.keys {
							for _, fkey := range childFieldInfo.keys {
								arrayFieldInfo.keys = append(arrayFieldInfo.keys, normalizeName(fmt.Sprintf("%s[%d].%s", akey, idx, fkey)))
							}
						}

						fieldsList = append(fieldsList, arrayFieldInfo)
					}
				}
			} else if arrayLength > 0 {
				// When the field is a slice/array of primitives, create a fieldInfo for each index
				for idx := 0; idx < arrayLength; idx++ {
					// copy index chain and append array index
					var cpy2 = make([]int, len(indexChain))
					copy(cpy2, indexChain)

					arrayFieldInfo := fieldInfo{
						IndexChain:   append(cpy2, idx),
						omitEmpty:    currFieldInfo.omitEmpty,
						defaultValue: currFieldInfo.defaultValue,
					}

					for _, akey := range currFieldInfo.keys {
						arrayFieldInfo.keys = append(arrayFieldInfo.keys, normalizeName(fmt.Sprintf("%s[%d]", akey, idx)))
					}

					fieldsList = append(fieldsList, arrayFieldInfo)
				}
			} else {
				fieldsList = append(fieldsList, currFieldInfo)
			}
		} else {
			fieldsList = append(fieldsList, currFieldInfo)
		}
	}
	return fieldsList
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

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

func isErrorType(outType reflect.Type) bool {
	if outType.Kind() != reflect.Interface {
		return false
	}

	return outType.Implements(errorInterface)
}
