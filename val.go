package gotemplconstr

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func marshalSimpleVal(val reflect.Value) (string, []byte, error) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(val.Uint(), 10), nil, nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, val.Type().Bits()), nil, nil
	case reflect.String:
		return val.String(), nil, nil
	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), nil, nil
	case reflect.Array:
		if typ.Elem().Kind() != reflect.Uint8 {
			break
		}
		// [...]byte
		var bytes []byte
		if val.CanAddr() {
			bytes = val.Slice(0, val.Len()).Bytes()
		} else {
			bytes = make([]byte, val.Len())
			reflect.Copy(reflect.ValueOf(bytes), val)
		}
		return "", bytes, nil
	case reflect.Slice:
		if typ.Elem().Kind() != reflect.Uint8 {
			break
		}
		// []byte
		return "", val.Bytes(), nil
	}

	return "", nil, fmt.Errorf("xml: unsupported type: %s", typ.String())
}

func getVal(key string, val reflect.Value) (reflect.Value, bool) {
	parts := strings.Split(key, ".")
	if len(parts) == 0 {
		return val, false
	}

	for _, part := range parts {
		val = val.FieldByName(part)
		if !val.IsValid() {
			return val, false
		}
	}

	return val, true
}
