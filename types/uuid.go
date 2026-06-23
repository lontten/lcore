package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	"strings"
)

type UUID uuid.UUID

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

func (u UUID) MarshalJSON() ([]byte, error) {
	all := strings.ReplaceAll(uuid.UUID(u).String(), "-", "")
	rs := []byte(fmt.Sprintf(`"%s"`, all))
	return rs, nil
}

func (u *UUID) UnmarshalJSON(src []byte) error {
	if len(src) != 34 {
		return errors.Errorf("invalid length for UUID: %v", len(src))
	}
	fromString, err := uuid.FromString(string(src[1 : len(src)-1]))
	if err != nil {
		return err
	}
	*u = UUID(fromString)
	return err
}

func (u UUID) Value() (driver.Value, error) {
	return uuid.UUID(u).String(), nil
}

func (u *UUID) Scan(v any) error {
	value, ok := v.(string)
	if ok {
		*u = UUID(uuid.FromStringOrNil(value))
		return nil
	}
	return fmt.Errorf("can not convert %v to uuid", v)
}

func Str2UUIDMust(v string) UUID {
	return UUID(uuid.FromStringOrNil(v))
}

func Str2UUIDMustP(v string) *UUID {
	fromString, err := uuid.FromString(v)
	if err != nil {
		return nil
	}
	u := UUID(fromString)
	return &u
}

func V4() uuid.UUID {
	v4, _ := uuid.NewV4()
	return v4
}

func V4p() *uuid.UUID {
	v4, _ := uuid.NewV4()
	return &v4
}

func NewV4() UUID {
	v4, _ := uuid.NewV4()
	return UUID(v4)
}

func NewV4P() *UUID {
	v4, err := uuid.NewV4()
	if err != nil {
		return nil
	}
	u := UUID(v4)
	return &u
}

func Str2UUID(v string) (UUID, error) {
	id, err := uuid.FromString(v)
	if err != nil {
		return UUID{}, err
	}
	return UUID(id), nil
}

type UUIDList []UUID

// Value 实现方法
func (p UUIDList) Value() (driver.Value, error) {
	var k []UUID
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
func (p *UUIDList) Scan(data any) error {
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
