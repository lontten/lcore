package types

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

// ArrayOf 由混合类型元素构造 PostgreSQL 数组字面量。
// 结构体元素若实现 driver.Valuer，则使用其 Value 输出；其余元素按类型格式化。
func ArrayOf(v ...any) Array {
	array := Array{}
	for _, i := range v {
		value := reflect.ValueOf(i)

		if value.Kind() == reflect.Struct {
			if valuer, ok := i.(driver.Valuer); ok {
				array.ints = append(array.ints, valuer)
				continue
			}
		}
		array.bases = append(array.bases, i)
	}
	return array
}

// Array 混合类型 PostgreSQL 数组的 GORM 自定义类型。
type Array struct {
	ints  []driver.Valuer
	bases []any
}

// Value 生成 PostgreSQL 数组文本字面量，如 {1,"hello"}。
func (a Array) Value() (driver.Value, error) {
	if len(a.ints) == 0 && len(a.bases) == 0 {
		return "{}", nil
	}
	str := "{"
	for _, e := range a.ints {
		value, err := e.Value()
		if err != nil {
			return nil, err
		}
		str += value.(string) + ","
	}
	for _, e := range a.bases {
		switch v := e.(type) {
		case int, int8, int16, int32, int64, float32, float64:
			str += fmt.Sprintf("%v", v) + ","
		default:
			str += fmt.Sprintf("\"%v\"", v) + ","
		}
	}
	str = str[:len(str)-1]
	str += "}"
	return str, nil
}
