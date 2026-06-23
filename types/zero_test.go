package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilToZero(t *testing.T) {
	as := assert.New(t)
	as.Equal(0, NilToZero[int](nil))
	as.Equal(42, NilToZero(ptr(42)))
}

func TestNilToZeroP(t *testing.T) {
	as := assert.New(t)
	p := NilToZeroP[int](nil)
	as.NotNil(p)
	as.Equal(0, *p)

	orig := ptr(7)
	got := NilToZeroP(orig)
	as.Equal(orig, got)
}

func ptr(v int) *int {
	return &v
}
