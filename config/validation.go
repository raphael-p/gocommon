package config

import (
	"fmt"
	"reflect"
	"regexp"
)

func traverseStructFields(
	reflectValue reflect.Value,
	jsonPath string,
	missingFields []string,
) []string {
	isJSONField := regexp.MustCompile(`^JSONField\[.+\]$`)
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Type().Field(i)
		tags := field.Tag

		jsonName := tags.Get("json")
		if jsonName == "" {
			jsonName = field.Name
		}

		if isJSONField.MatchString(field.Type.Name()) {
			jsonField := reflectValue.FieldByName(field.Name)
			isOptional := tags.Get("optional") == "true"
			isNullable := tags.Get("nullable") == "true"
			isZeroable := isOptional || isNullable || tags.Get("zeroable") == "true"
			isSet := jsonField.FieldByName("IsSet").Bool()
			isNull := jsonField.FieldByName("IsNull").Bool()
			value := jsonField.FieldByName("Value")
			isZero := value.IsZero()
			isStruct := value.Type().Kind() == reflect.Struct
			if (!isSet && !isOptional) ||
				(isNull && !isNullable) ||
				(isZero && !isZeroable && !isStruct) {
				missingFields = append(missingFields, jsonPath+jsonName)
			}
			if isStruct {
				missingFields = traverseStructFields(
					value,
					jsonPath+jsonName+".",
					missingFields,
				)
			}
		} else if field.Type.Kind() == reflect.Struct {
			missingFields = traverseStructFields(
				reflectValue.FieldByName(field.Name),
				jsonPath+jsonName+".",
				missingFields,
			)
		} else {
			if reflectValue.FieldByName(field.Name).IsZero() {
				missingFields = append(missingFields, jsonPath+jsonName)
			}
		}
	}
	return missingFields
}

// Finds any fields from the struct which are zero-valued.
// `value` must be a `struct` or its pointer.
// Accepted tags:
//
//	`json:"nameInJson"` -> maps JSON to struct fields when deserialising JSON
//	`optional:"true"` -> for JSONField only, allows field to not be set
//	`nullable:"true"` -> for JSONField only, allows field to be set to null
//	`zeroable:"true"` -> for JSONField only, allows field to be set to zero-value
//
// Used to validate the deserialisation of a JSON document.
func structFromJSON(value any) ([]string, error) {
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	if reflectValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected `value` to be a struct or its pointer, got %T", value)
	}
	return traverseStructFields(reflectValue, "", []string{}), nil
}
