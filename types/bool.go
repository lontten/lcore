package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

// BoolList 映射 PostgreSQL bool[]，实现 driver.Valuer、sql.Scanner 及 sort.Interface。
type BoolList []bool

// Value 返回 bool[] 的 PostgreSQL 文本字面量，如 {true,false}。
func (p BoolList) Value() (driver.Value, error) {
	marshal, err := json.Marshal([]bool(p))
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

// Scan 从 PostgreSQL bool[] 文本字面量或驱动返回值解析；nil 时不修改接收方。
func (p *BoolList) Scan(data any) error {
	if data == nil {
		return nil
	}
	var list []bool
	if err := scanPgArray(pgtype.BoolArrayOID, data, &list); err != nil {
		return err
	}
	*p = BoolList(list)
	return nil
}

// Len 返回元素个数，实现 sort.Interface。
func (p BoolList) Len() int {
	return len(p)
}

// Less 定义排序：false 排在 true 之前。
func (p BoolList) Less(i, j int) bool {
	return !p[i] && p[j]
}

// Swap 交换两个元素，实现 sort.Interface。
func (p BoolList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
