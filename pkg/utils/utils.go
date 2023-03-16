package utils

import (
	"reflect"
)

// StructToMap converts a struct to a map, where the keys of the map are the
// field names of the struct and the values of the map are the corresponding
// field values of the struct.
func StructToMap(obj interface{}) map[string]interface{} {
	// create an empty map to hold the result
	result := make(map[string]interface{})

	// get the reflect value of the object
	val := reflect.ValueOf(obj)

	// get the reflect type of the object
	typ := reflect.TypeOf(obj)

	// iterate over the fields of the struct
	for i := 0; i < val.NumField(); i++ {
		// get the reflect value of the field
		field := val.Field(i)

		// get the name of the field from the struct's type
		fieldName := typ.Field(i).Name

		// set the field's value in the map
		result[fieldName] = field.Interface()
	}

	return result
}
