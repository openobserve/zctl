package utils

import (
	"fmt"
	"reflect"
)

func StructToMap(s interface{}) (map[string]interface{}, error) {
	return structToMapRec(s, reflect.TypeOf(s), reflect.ValueOf(s))
}

func structToMapRec(s interface{}, typ reflect.Type, val reflect.Value) (map[string]interface{}, error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input should be a struct, got %s", val.Kind())
	}

	result := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {

		fieldName := typ.Field(i).Tag.Get("yaml")
		fieldValue := val.Field(i)

		if fieldValue.Kind() == reflect.Struct {
			nestedMap, err := structToMapRec(fieldValue.Interface(), fieldValue.Type(), fieldValue)
			if err != nil {
				return nil, err
			}
			result[fieldName] = nestedMap
		} else {
			result[fieldName] = fieldValue.Interface()
		}
	}

	return result, nil
}
