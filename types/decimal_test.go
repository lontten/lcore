package types

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecimalList_ValueScan(t *testing.T) {
	as := assert.New(t)
	list := DecimalList{
		decimal.NewFromInt(1),
		decimal.RequireFromString("2.50"),
	}

	value, err := list.Value()
	as.NoError(err)

	var scanned DecimalList
	as.NoError(scanned.Scan(value))
	as.Equal(2, len(scanned))
	as.True(decimal.NewFromInt(1).Equal(scanned[0]))
	as.True(decimal.RequireFromString("2.50").Equal(scanned[1]))
}

func TestDecimalList_ScanInvalid(t *testing.T) {
	as := assert.New(t)
	var list DecimalList
	as.Error(list.Scan(`{"not-a-number"}`))
}

func TestToDecimal(t *testing.T) {
	as := assert.New(t)
	as.True(decimal.NewFromInt(42).Equal(ToDecimal(42)))
	as.True(decimal.NewFromInt(42).Equal(ToDecimal(int64(42))))
	as.True(decimal.NewFromFloat(3.14).Equal(ToDecimal(3.14)))
	as.True(decimal.RequireFromString("1.23").Equal(ToDecimal("1.23")))
	as.Panics(func() { ToDecimal("bad") })
	as.Panics(func() { ToDecimal(struct{}{}) })
}

func TestDecimalConverters(t *testing.T) {
	as := assert.New(t)
	as.True(decimal.NewFromInt(10).Equal(IntToDecimal(10)))
	as.True(decimal.NewFromInt(20).Equal(Int64ToDecimal(20)))
	as.True(decimal.NewFromFloat(1.5).Equal(Float64ToDecimal(1.5)))
	as.True(decimal.RequireFromString("9.99").Equal(StringToDecimal("9.99")))
	as.Panics(func() { StringToDecimal("x") })
}

func TestDecimalList_Empty(t *testing.T) {
	as := assert.New(t)
	empty := DecimalList{}
	val, err := empty.Value()
	as.NoError(err)
	require.NotNil(t, val)

	var scanned DecimalList
	as.NoError(scanned.Scan(val))
	as.Empty(scanned)
}
