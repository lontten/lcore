package structutil

import "reflect"

func Struct2StringMap(s any) map[string]string {
	val := reflect.ValueOf(s)
	relType := val.Type()
	numField := relType.NumField()
	m := make(map[string]string, numField)
	for i := 0; i < numField; i++ {
		m[relType.Field(i).Name] = val.Field(i).String()
	}
	return m
}
