package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNullUint64_IsZero(t *testing.T) {
	as := assert.New(t)
	as.True(NullUint64{}.IsZero())
	as.True(NullUint64{Valid: true, Uint64: 0}.IsZero())
	as.False(NewNullUint64(42).IsZero())
}

func TestNullUint64_Scan(t *testing.T) {
	as := assert.New(t)
	var n NullUint64

	as.NoError(n.Scan(nil))
	as.False(n.Valid)

	as.NoError(n.Scan(int64(100)))
	as.True(n.Valid)
	as.Equal(uint64(100), n.Uint64)

	as.NoError(n.Scan(uint64(200)))
	as.Equal(uint64(200), n.Uint64)

	as.NoError(n.Scan("300"))
	as.Equal(uint64(300), n.Uint64)

	as.NoError(n.Scan([]byte("400")))
	as.Equal(uint64(400), n.Uint64)

	as.Error(n.Scan(int64(-1)))
	as.Error(n.Scan(struct{}{}))
}

func TestNullUint64_Value(t *testing.T) {
	as := assert.New(t)
	n := NewNullUint64(42)
	val, err := n.Value()
	as.NoError(err)
	as.Equal(int64(42), val)

	invalid := NullUint64{}
	val, err = invalid.Value()
	as.NoError(err)
	as.Nil(val)
}

func TestNullUint64_JSON(t *testing.T) {
	as := assert.New(t)
	n := NewNullUint64(99)
	data, err := n.MarshalJSON()
	as.NoError(err)
	as.Equal("99", string(data))

	var dst NullUint64
	as.NoError(json.Unmarshal(data, &dst))
	as.True(dst.Valid)
	as.Equal(uint64(99), dst.Uint64)

	as.NoError(json.Unmarshal([]byte("null"), &dst))
	as.False(dst.Valid)
}

func TestNullUint64_String(t *testing.T) {
	as := assert.New(t)
	as.Equal("NULL", NullUint64{}.String())
	as.Equal("42", NewNullUint64(42).String())
}

func TestNewNullUint64(t *testing.T) {
	as := assert.New(t)
	n := NewNullUint64(7)
	require.True(t, n.Valid)
	as.Equal(uint64(7), n.Uint64)
}
