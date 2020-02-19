package utils

import (
	"errors"
	"fmt"
	"reflect"
)

// SetField 给定 k,v, 填充struct中指定字段的值
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

// FillStruct 使用map的值填充struct
func FillStruct(obj interface{}, m map[string]interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			continue
		}
		mval, ok := m[tag]
		if !ok {
			continue
		}
		val := reflect.ValueOf(mval)

		if field.Type != val.Type() {
			return errors.New("Provided value type didn't match obj field type")
		}
		structValue.Field(i).Set(val)
	}
	return nil
}
