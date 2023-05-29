package utils

import (
	"reflect"
	"unsafe"
)

func StructMapTagValueToFieldValue(src any, tagName string, includeWithoutTag, includeEmptyTag bool) (m map[string]any) {
	m = map[string]any{}
	StructFieldsTagsIterator(src, tagName, func(fieldName, tagValue string, hasTag bool, vType reflect.Type, value reflect.Value) {
		if !hasTag && includeWithoutTag {
			hasTag = true
			tagValue = fieldName
		}
		if hasTag && (includeEmptyTag || tagValue != "") {
			m[tagValue] = value.Interface()
		}
	})
	return
}

func StructFieldsTagsIterator(src any, tagName string, f func(fieldName, tagValue string, hasTag bool, vType reflect.Type, value reflect.Value)) {
	val := ReflectDereference(src)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := ReflectDereference(val.Field(i))

		if fieldType.Type.Kind() == reflect.Struct {
			StructFieldsTagsIterator(fieldValue, tagName, f)
		} else {
			tagValue, hasTagValue := typ.Field(i).Tag.Lookup(tagName)
			if fieldType.IsExported() {
				f(fieldType.Name, tagValue, hasTagValue, fieldType.Type, fieldValue)
			}
		}
	}
}

func ReflectDereference(src any) (v reflect.Value) {
	var val reflect.Value
	var ok bool
	if val, ok = src.(reflect.Value); !ok {
		val = reflect.ValueOf(src)
	}
	//println("DEREF:", val.Kind().String(), val.Type().String())
	if val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface {
		return ReflectDereference(val.Elem())
	}
	return val
}

// ReflectGetFieldByName retrieves even private fields from the provided structure.
// You can then use Interface() on them.
func ReflectGetFieldByName[T any](from *T, fieldName string) reflect.Value {
	value := reflect.ValueOf(from).Elem().FieldByName(fieldName)
	value = reflect.NewAt(value.Type(), unsafe.Pointer(value.UnsafeAddr())).Elem()
	return value
}
