package paginator

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
	"time"
)

/* public */

// Encode encodes properties in order defined by keys on the struct of v
func Encode(v reflect.Value, keys []string) string {
	fields := make([]string, len(keys))
	for index, key := range keys {
		if v.Kind() == reflect.Ptr {
			fields[index] = convert(reflect.Indirect(v).FieldByName(key).Interface())
		} else {
			fields[index] = convert(v.FieldByName(key).Interface())
		}
	}
	return encodeBase64(fields)
}

// Decode decodes cursor into values in the same order as encoding
func Decode(cursor string) []interface{} {
	// @TODO: backward-compatibility
	fieldsWithType := decodeBase64(cursor)
	fields := make([]interface{}, len(fieldsWithType))
	for index, fieldWithType := range fieldsWithType {
		fields[index] = revert(fieldWithType)
	}
	return fields
}

/* private */

type fieldType string

const (
	fieldString fieldType = "STRING"
	fieldTime   fieldType = "TIME"
)

func convert(field interface{}) string {
	switch field.(type) {
	case time.Time:
		return fmt.Sprintf("%s?%s", field.(time.Time).UTC().Format(time.RFC3339Nano), fieldTime)
	default:
		return fmt.Sprintf("%v?%s", field, fieldString)
	}
}

func revert(fieldWithType string) interface{} {
	field, fieldType := parseFieldWithType(fieldWithType)
	switch fieldType {
	case fieldTime:
		t, err := time.Parse(time.RFC3339Nano, field)
		if err != nil {
			t = time.Now().UTC()
		}
		return t
	default:
		return field
	}
}

func parseFieldWithType(fieldWithType string) (string, fieldType) {
	sep := strings.LastIndex(fieldWithType, "?")
	field := fieldWithType[:sep]
	fieldType := fieldType(fieldWithType[sep+1:])
	return field, fieldType
}

func reverse(v reflect.Value) reflect.Value {
	result := reflect.MakeSlice(v.Type(), 0, v.Cap())
	for i := v.Len() - 1; i >= 0; i-- {
		result = reflect.Append(result, v.Index(i))
	}
	return result
}

func encodeBase64(fields []string) string {
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(fields, ",")))
}

func decodeBase64(cursor string) []string {
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil
	}
	return strings.Split(string(b), ",")
}
