package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHelpers(t *testing.T) {
	as := assert.New(t)

	i := NewInt(42)
	as.NotNil(i)
	as.Equal(42, *i)
	*i = 99
	as.Equal(99, *i)

	i8 := NewInt8(8)
	as.Equal(int8(8), *i8)

	i16 := NewInt16(16)
	as.Equal(int16(16), *i16)

	i32 := NewInt32(32)
	as.Equal(int32(32), *i32)

	i64 := NewInt64(64)
	as.Equal(int64(64), *i64)

	s := NewString("hello")
	as.Equal("hello", *s)

	b := NewBool(true)
	as.True(*b)
}
