package types

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestScanPgArray_Nil(t *testing.T) {
	as := assert.New(t)
	var dst []string
	as.NoError(scanPgArray(pgtype.TextArrayOID, nil, &dst))
	as.Nil(dst)
}

func TestScanPgArray_StringBoolArray(t *testing.T) {
	as := assert.New(t)
	var dst []bool
	as.NoError(scanPgArray(pgtype.BoolArrayOID, "{true,false}", &dst))
	as.Equal([]bool{true, false}, dst)
}

func TestScanPgArray_Bytes(t *testing.T) {
	as := assert.New(t)
	var list BoolList
	as.NoError(scanPgArray(pgtype.BoolArrayOID, []byte("{true}"), &list))
	as.Equal(BoolList{true}, list)
}

func TestScanPgArray_TextArray(t *testing.T) {
	as := assert.New(t)
	var dst []string
	as.NoError(scanPgArray(pgtype.TextArrayOID, `{"a","b"}`, &dst))
	as.Equal([]string{"a", "b"}, dst)
}

func TestScanPgArray_InvalidType(t *testing.T) {
	as := assert.New(t)
	var dst []string
	as.Error(scanPgArray(pgtype.TextArrayOID, 123, &dst))
}
