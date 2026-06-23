package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNum2Str(t *testing.T) {
	as := assert.New(t)
	as.Equal("42", Num2Str(42))
	as.Equal("-1", Num2Str(int8(-1)))

	p := Num2StrP(10)
	require.NotNil(t, p)
	as.Equal("10", *p)
}

func TestStr2Num_Signed(t *testing.T) {
	as := assert.New(t)

	v, err := Str2Num[int]("42")
	as.NoError(err)
	as.Equal(42, v)

	v8, err := Str2Num[int8]("0x10")
	as.NoError(err)
	as.Equal(int8(16), v8)

	_, err = Str2Num[int8]("999")
	as.Error(err)

	_, err = Str2Num[int]("not-a-number")
	as.Error(err)
}

func TestStr2Num_Unsigned(t *testing.T) {
	as := assert.New(t)
	v, err := Str2Num[uint32]("100")
	as.NoError(err)
	as.Equal(uint32(100), v)

	_, err = Str2Num[uint8]("-1")
	as.Error(err)
}

func TestStr2NumMust(t *testing.T) {
	as := assert.New(t)
	as.Equal(42, Str2NumMust[int]("42"))
	as.Panics(func() { Str2NumMust[int]("bad") })

	p := Str2NumMustP[int32]("10")
	as.NotNil(p)
	as.Equal(int32(10), *p)
}
