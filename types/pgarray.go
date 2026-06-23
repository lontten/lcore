package types

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

var pgMap = pgtype.NewMap()

// scanPgArray 将 PostgreSQL 数组文本字面量解析到 dst（通常为 FlatArray 或 []T）。
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
