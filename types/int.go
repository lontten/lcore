package types

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

// IntList 映射 PostgreSQL 整数数组（以 text[] 形式读写），实现 driver.Valuer 与 sql.Scanner。
type IntList []int

// Value 返回整数数组的 PostgreSQL 文本字面量，如 {1,2,3}。
func (p IntList) Value() (driver.Value, error) {
	marshal, err := json.Marshal([]int(p))
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

// Scan 从 PostgreSQL 数组文本解析为 []int；nil 时不修改接收方。
func (p *IntList) Scan(data any) error {
	if data == nil {
		return nil
	}
	var list []string
	if err := scanPgArray(pgtype.TextArrayOID, data, &list); err != nil {
		return err
	}
	ints := make(IntList, len(list))
	for i, s := range list {
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		ints[i] = n
	}
	*p = ints
	return nil
}
