package types

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testUUID = "550e8400-e29b-41d4-a716-446655440000"

func TestUUID_StringAndDB(t *testing.T) {
	as := assert.New(t)
	u, err := Str2UUID(testUUID)
	as.NoError(err)
	as.Equal(testUUID, u.String())

	val, err := u.Value()
	as.NoError(err)
	as.Equal(testUUID, val)

	var scanned UUID
	as.NoError(scanned.Scan(testUUID))
	as.Equal(u, scanned)

	as.Error(scanned.Scan(123))
}

func TestUUID_JSON(t *testing.T) {
	as := assert.New(t)
	u, err := Str2UUID(testUUID)
	as.NoError(err)

	data, err := u.MarshalJSON()
	as.NoError(err)
	as.Equal(`"550e8400e29b41d4a716446655440000"`, string(data))

	var dst UUID
	as.NoError(dst.UnmarshalJSON(data))
	as.Equal(u, dst)

	as.Error(dst.UnmarshalJSON([]byte(`"tooshort"`)))
	as.Error(dst.UnmarshalJSON([]byte(`not-json`)))
}

func TestStr2UUIDMust(t *testing.T) {
	as := assert.New(t)
	u := Str2UUIDMust(testUUID)
	as.Equal(testUUID, u.String())
	as.Equal(UUID{}, Str2UUIDMust("bad"))

	p := Str2UUIDMustP(testUUID)
	require.NotNil(t, p)
	as.Equal(testUUID, p.String())
	as.Nil(Str2UUIDMustP("bad"))
}

func TestNewV4(t *testing.T) {
	as := assert.New(t)
	u := NewV4()
	as.NotEqual(UUID{}, u)
	as.Len(u.String(), 36)

	p := NewV4P()
	require.NotNil(t, p)

	raw := V4()
	as.NotEqual(uuid.Nil, raw)
	rawP := V4p()
	require.NotNil(t, rawP)
}

func TestUUIDList_ValueScan(t *testing.T) {
	as := assert.New(t)
	u1, _ := Str2UUID(testUUID)
	u2 := NewV4()
	list := UUIDList{u1, u2}

	value, err := list.Value()
	as.NoError(err)

	var scanned UUIDList
	as.NoError(scanned.Scan(value))
	as.Equal(2, len(scanned))
	as.Equal(u1, scanned[0])
	as.Equal(u2, scanned[1])
}

func TestUUIDList_Empty(t *testing.T) {
	as := assert.New(t)
	empty := UUIDList{}
	val, err := empty.Value()
	as.NoError(err)

	var scanned UUIDList
	as.NoError(scanned.Scan(val))
	as.Empty(scanned)
}
