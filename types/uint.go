package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

type NullUint64 struct {
	Uint64 uint64
	Valid  bool
}

// NewNullUint64 创建有效的 NullUint64
func NewNullUint64(v uint64) NullUint64 {
	return NullUint64{Uint64: v, Valid: true}
}

// IsZero 判断是否是零值或无效值
func (n NullUint64) IsZero() bool {
	return !n.Valid || n.Uint64 == 0
}

// Scan 支持更多类型的扫描
func (n *NullUint64) Scan(value interface{}) error {
	if value == nil {
		n.Uint64, n.Valid = 0, false
		return nil
	}

	n.Valid = true

	switch v := value.(type) {
	case int64:
		if v < 0 {
			return fmt.Errorf("uint64 cannot be negative: %d", v)
		}
		n.Uint64 = uint64(v)
	case int32:
		n.Uint64 = uint64(v)
	case int:
		if v < 0 {
			return fmt.Errorf("uint64 cannot be negative: %d", v)
		}
		n.Uint64 = uint64(v)
	case uint64:
		n.Uint64 = v
	case uint32:
		n.Uint64 = uint64(v)
	case uint:
		n.Uint64 = uint64(v)
	case []byte:
		// 从字符串解析
		return n.parseString(string(v))
	case string:
		return n.parseString(v)
	default:
		return fmt.Errorf("unsupported Scan type for NullUint64: %T", value)
	}

	return nil
}

func (n *NullUint64) parseString(s string) error {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	n.Uint64 = v
	n.Valid = true
	return nil
}

// Value 返回数据库值
func (n NullUint64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	// 注意：返回 int64，大部分数据库驱动期望这样
	return int64(n.Uint64), nil
}

// JSON 序列化/反序列化
func (n NullUint64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Uint64)
}

func (n *NullUint64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}

	var v uint64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	n.Uint64 = v
	n.Valid = true
	return nil
}

// String 方法
func (n NullUint64) String() string {
	if !n.Valid {
		return "NULL"
	}
	return strconv.FormatUint(n.Uint64, 10)
}
