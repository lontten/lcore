package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

// DecimalList 映射 PostgreSQL numeric/text[]（以 text[] 读写 decimal 字符串），供 GORM 自定义字段使用。
type DecimalList []decimal.Decimal

// Value 返回 decimal 数组的 PostgreSQL 文本字面量。
func (p DecimalList) Value() (driver.Value, error) {
	marshal, err := json.Marshal([]decimal.Decimal(p))
	if err != nil {
		return nil, err
	}
	s := string(marshal)
	if s != "null" {
		s = "{" + s[1:len(s)-1] + "}"
	} else {
		s = "{}"
	}
	return s, nil
}

// Scan 从 PostgreSQL text[] 解析为 []decimal.Decimal；nil 时不修改接收方。
func (p *DecimalList) Scan(data any) error {
	if data == nil {
		return nil
	}
	var list []string
	if err := scanPgArray(pgtype.TextArrayOID, data, &list); err != nil {
		return err
	}
	decimals := make([]decimal.Decimal, len(list))
	for i, s := range list {
		fromString, err := decimal.NewFromString(s)
		if err != nil {
			return err
		}
		decimals[i] = fromString
	}
	*p = decimals
	return nil
}

// ToDecimal 将整数、浮点数或字符串转换为 decimal.Decimal。
// 支持的类型：int、int8~int64、uint、uint8~uint64、float32、float64、string。
// 字符串解析失败或类型不支持时 panic。
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

// IntToDecimal 将 int 转为 decimal.Decimal。
func IntToDecimal(value int) decimal.Decimal {
	return decimal.NewFromInt(int64(value))
}

// Int64ToDecimal 将 int64 转为 decimal.Decimal。
func Int64ToDecimal(value int64) decimal.Decimal {
	return decimal.NewFromInt(value)
}

// Float64ToDecimal 将 float64 转为 decimal.Decimal。
func Float64ToDecimal(value float64) decimal.Decimal {
	return decimal.NewFromFloat(value)
}

// StringToDecimal 将字符串转为 decimal.Decimal，解析失败时 panic。
func StringToDecimal(value string) decimal.Decimal {
	result, err := decimal.NewFromString(value)
	if err != nil {
		panic(err)
	}
	return result
}
