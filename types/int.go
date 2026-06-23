package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

type IntList []int

func (p IntList) Value() (driver.Value, error) {
	marshal, err := json.Marshal([]int(p))
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

func (p *IntList) Scan(data any) error {
	var list []string
	if err := scanPgArray(pgtype.TextArrayOID, data, &list); err != nil {
		return err
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return json.Unmarshal(marshal, p)
}
