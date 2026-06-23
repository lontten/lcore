package types

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

var pgMap = pgtype.NewMap()

func scanPgArray(oid uint32, data any, dst any) error {
	if data == nil {
		return nil
	}
	var src []byte
	switch v := data.(type) {
	case string:
		src = []byte(v)
	case []byte:
		src = v
	default:
		return fmt.Errorf("cannot scan %T into pg array", data)
	}
	return pgMap.Scan(oid, pgtype.TextFormatCode, src, dst)
}
