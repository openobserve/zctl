package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"time"
)

func StructToMap2(s interface{}) (map[string]interface{}, error) {
	return structToMapRec2(reflect.TypeOf(s), reflect.ValueOf(s))
}

func structToMapRec2(typ reflect.Type, val reflect.Value) (map[string]interface{}, error) {
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

		switch fieldValue.Kind() {
		case reflect.Struct:
			nestedMap, err := structToMapRec2(fieldValue.Type(), fieldValue)
			if err != nil {
				return nil, err
			}
			result[fieldName] = nestedMap
		case reflect.Slice, reflect.Array:
			length := fieldValue.Len()
			slice := make([]interface{}, length)
			for j := 0; j < length; j++ {
				elem := fieldValue.Index(j)
				if elem.Kind() == reflect.Struct {
					nestedMap, err := structToMapRec2(elem.Type(), elem)
					if err != nil {
						return nil, err
					}
					slice[j] = nestedMap
				} else {
					slice[j] = elem.Interface()
				}
			}
			result[fieldName] = slice
		default:
			result[fieldName] = fieldValue.Interface()
		}
	}

	return result, nil
}

// Oretty print go struct
func PrettyPrint(s1 map[string]interface{}) {
	b, err := json.MarshalIndent(s1, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
}

func GenerateReleaseIdentifier() string {
	// Set the seed value for the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 99999
	randomNumber := rand.Intn(100000)

	// If the number is less than 10000, add leading zeros
	if randomNumber < 10000 {
		randomNumber += 10000
	}

	return fmt.Sprintf("%d", randomNumber)
}
