package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

// BoolList PostgreSQL bool[] 的 GORM 自定义类型。
type BoolList []bool

func (p BoolList) Value() (driver.Value, error) {
	marshal, err := json.Marshal([]bool(p))
	if err != nil {
		return nil, err
	}
	s := string(marshal)
	if s != "null" {
		s = s[:0] + "{" + s[1:len(s)-1] + "}" + s[len(s):]
	} else {
		s = "{}"
	}
	return s, nil
}

func (p *BoolList) Scan(data any) error {
	var list []bool
	if err := scanPgArray(pgtype.BoolArrayOID, data, &list); err != nil {
		return err
	}
	*p = BoolList(list)
	return nil
}

func (p BoolList) Len() int {
	return len(p)
}

func (p BoolList) Less(i, j int) bool {
	return !p[i] && p[j]
}

func (p BoolList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
