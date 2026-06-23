package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
)

// UUID 是 github.com/gofrs/uuid.UUID 的类型别名，提供 JSON/数据库互操作。
type UUID uuid.UUID

// String 返回带连字符的标准 UUID 字符串（小写）。
func (u UUID) String() string {
	return uuid.UUID(u).String()
}

// MarshalJSON 输出无连字符的 32 位十六进制字符串，如 "550e8400e29b41d4a716446655440000"。
func (u UUID) MarshalJSON() ([]byte, error) {
	all := strings.ReplaceAll(uuid.UUID(u).String(), "-", "")
	return []byte(fmt.Sprintf(`"%s"`, all)), nil
}

// UnmarshalJSON 解析无连字符的 32 位十六进制 JSON 字符串（总长 34 字节含引号）。
// 长度不符或格式非法时返回错误。
func (u *UUID) UnmarshalJSON(src []byte) error {
	if len(src) != 34 {
		return errors.Errorf("invalid length for UUID: %v", len(src))
	}
	fromString, err := uuid.FromString(string(src[1 : len(src)-1]))
	if err != nil {
		return err
	}
	*u = UUID(fromString)
	return nil
}

// Value 写入数据库的标准 UUID 字符串（带连字符）。
func (u UUID) Value() (driver.Value, error) {
	return uuid.UUID(u).String(), nil
}

// Scan 从数据库字符串读取 UUID；非法字符串写入 nil UUID 而不报错。
func (u *UUID) Scan(v any) error {
	value, ok := v.(string)
	if ok {
		*u = UUID(uuid.FromStringOrNil(value))
		return nil
	}
	return fmt.Errorf("can not convert %v to uuid", v)
}

// Str2UUIDMust 解析 UUID 字符串，失败时返回 nil UUID（不 panic）。
func Str2UUIDMust(v string) UUID {
	return UUID(uuid.FromStringOrNil(v))
}

// Str2UUIDMustP 解析 UUID 字符串，成功返回指针，失败返回 nil。
func Str2UUIDMustP(v string) *UUID {
	fromString, err := uuid.FromString(v)
	if err != nil {
		return nil
	}
	u := UUID(fromString)
	return &u
}

// V4 生成随机 UUID（gofrs/uuid 类型）。
func V4() uuid.UUID {
	v4, _ := uuid.NewV4()
	return v4
}

// V4p 生成随机 UUID 指针（gofrs/uuid 类型）。
func V4p() *uuid.UUID {
	v4, _ := uuid.NewV4()
	return &v4
}

// NewV4 生成随机 UUID（types.UUID 类型）。
func NewV4() UUID {
	v4, _ := uuid.NewV4()
	return UUID(v4)
}

// NewV4P 生成随机 UUID 指针（types.UUID 类型），失败时返回 nil。
func NewV4P() *UUID {
	v4, err := uuid.NewV4()
	if err != nil {
		return nil
	}
	u := UUID(v4)
	return &u
}

// Str2UUID 解析 UUID 字符串，失败返回 error。
func Str2UUID(v string) (UUID, error) {
	id, err := uuid.FromString(v)
	if err != nil {
		return UUID{}, err
	}
	return UUID(id), nil
}

// UUIDList 映射 PostgreSQL uuid[]，实现 driver.Valuer 与 sql.Scanner。
type UUIDList []UUID

// Value 返回 uuid[] 的 PostgreSQL 文本字面量。
func (p UUIDList) Value() (driver.Value, error) {
	marshal, err := json.Marshal([]UUID(p))
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

// Scan 从 PostgreSQL uuid[] 解析；nil 时不修改接收方。
func (p *UUIDList) Scan(data any) error {
	if data == nil {
		return nil
	}
	var uuids pgtype.FlatArray[pgtype.UUID]
	if err := scanPgArray(pgtype.UUIDArrayOID, data, &uuids); err != nil {
		return err
	}
	list := make([]UUID, len(uuids))
	for i, element := range uuids {
		if element.Valid {
			u, err := uuid.FromBytes(element.Bytes[:])
			if err != nil {
				return err
			}
			list[i] = UUID(u)
		}
	}
	*p = list
	return nil
}
