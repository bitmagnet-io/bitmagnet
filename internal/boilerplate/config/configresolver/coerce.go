package configresolver

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var durationType = reflect.TypeOf(time.Duration(0))

func coerceStringValue(stringValue string, valueType reflect.Type) (interface{}, error) {
	// todo Fill this out
	switch valueType {
	case durationType:
		return time.ParseDuration(stringValue)
	}
	switch valueType.Kind() {
	case reflect.String:
		return stringValue, nil
	case reflect.Bool:
		switch strings.ToLower(stringValue) {
		case "true", "1":
			return true, nil
		case "false", "0":
			return false, nil
		}
	case reflect.Int, reflect.Int64:
		return strconv.Atoi(stringValue)
	case reflect.Uint, reflect.Uint16, reflect.Uint64:
		return strconv.ParseUint(stringValue, 10, 64)
	case reflect.Slice:
		strValues := strings.Split(stringValue, ",")
		values := make([]interface{}, len(strValues))
		for i, strValue := range strValues {
			coercedValue, err := coerceStringValue(strValue, valueType.Elem())
			if err != nil {
				return nil, err
			}
			values[i] = coercedValue
		}
		return values, nil
	}
	return nil, fmt.Errorf("cannot coerce value to type %v", valueType)
}
