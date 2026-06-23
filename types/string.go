package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

type StringList []string

func (p StringList) Value() (driver.Value, error) {
	marshal, err := json.Marshal([]string(p))
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

func (p *StringList) Scan(data any) error {
	var list []string
	if err := scanPgArray(pgtype.TextArrayOID, data, &list); err != nil {
		return err
	}
	*p = StringList(list)
	return nil
}

func (p StringList) Len() int {
	return len(p)
}

func (p StringList) Less(i, j int) bool {
	return p[i] < p[j]
}

func (p StringList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
