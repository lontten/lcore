package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

type DecimalList []decimal.Decimal

// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p DecimalList) Value() (driver.Value, error) {
	var k []decimal.Decimal
	k = p
	marshal, err := json.Marshal(k)
	if err != nil {
		return nil, err
	}
	var s = string(marshal)
	if s != "null" {
		s = s[:0] + "{" + s[1:len(s)-1] + "}" + s[len(s):]
	} else {
		s = "{}"
	}
	return s, nil
}

// Scan 实现方法
func (p *DecimalList) Scan(data any) error {
	array := pgtype.VarcharArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []decimal.Decimal
	list = make([]decimal.Decimal, len(array.Elements))
	for i, element := range array.Elements {
		fromString, err := decimal.NewFromString(element.String)
		if err != nil {
			return err
		}
		list[i] = fromString
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}

// ToDecimal 将整数、浮点数或字符串转换为 decimal.Decimal
// 支持的类型: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string
// 如果转换失败会直接 panic
func ToDecimal(value any) decimal.Decimal {
	switch v := value.(type) {
	case int:
		return decimal.NewFromInt(int64(v))
	case int8:
		return decimal.NewFromInt(int64(v))
	case int16:
		return decimal.NewFromInt(int64(v))
	case int32:
		return decimal.NewFromInt(int64(v))
	case int64:
		return decimal.NewFromInt(v)
	case uint:
		return decimal.NewFromInt(int64(v))
	case uint8:
		return decimal.NewFromInt(int64(v))
	case uint16:
		return decimal.NewFromInt(int64(v))
	case uint32:
		return decimal.NewFromInt(int64(v))
	case uint64:
		return decimal.NewFromInt(int64(v))
	case float32:
		return decimal.NewFromFloat32(v)
	case float64:
		return decimal.NewFromFloat(v)
	case string:
		result, err := decimal.NewFromString(v)
		if err != nil {
			panic(err)
		}
		return result
	default:
		panic("unsupported type for decimal conversion")
	}
}

// IntToDecimal 将整数转换为 decimal.Decimal
func IntToDecimal(value int) decimal.Decimal {
	return decimal.NewFromInt(int64(value))
}

// Int64ToDecimal 将 int64 转换为 decimal.Decimal
func Int64ToDecimal(value int64) decimal.Decimal {
	return decimal.NewFromInt(value)
}

// Float64ToDecimal 将 float64 转换为 decimal.Decimal
func Float64ToDecimal(value float64) decimal.Decimal {
	return decimal.NewFromFloat(value)
}

// StringToDecimal 将字符串转换为 decimal.Decimal
// 如果转换失败会直接 panic
func StringToDecimal(value string) decimal.Decimal {
	result, err := decimal.NewFromString(value)
	if err != nil {
		panic(err)
	}
	return result
}
