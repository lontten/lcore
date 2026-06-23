package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

// StringList 映射 PostgreSQL text[]，实现 driver.Valuer、sql.Scanner 及 sort.Interface。
type StringList []string

// Value 返回 text[] 的 PostgreSQL 文本字面量，如 {"a","b"}。
func (p StringList) Value() (driver.Value, error) {
	marshal, err := json.Marshal([]string(p))
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

// Scan 从 PostgreSQL text[] 解析；nil 时不修改接收方。
func (p *StringList) Scan(data any) error {
	if data == nil {
		return nil
	}
	var list []string
	if err := scanPgArray(pgtype.TextArrayOID, data, &list); err != nil {
		return err
	}
	*p = StringList(list)
	return nil
}

// Len 返回元素个数，实现 sort.Interface。
func (p StringList) Len() int {
	return len(p)
}

// Less 按字典序比较，实现 sort.Interface。
func (p StringList) Less(i, j int) bool {
	return p[i] < p[j]
}

// Swap 交换两个元素，实现 sort.Interface。
func (p StringList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
