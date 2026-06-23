package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayOf_Empty(t *testing.T) {
	as := assert.New(t)
	val, err := ArrayOf().Value()
	as.NoError(err)
	as.Equal("{}", val)
}

func TestArrayOf_Integers(t *testing.T) {
	as := assert.New(t)
	arr := ArrayOf(1, 2, 3)
	val, err := arr.Value()
	as.NoError(err)
	as.Equal("{1,2,3}", val)
}

func TestArrayOf_Strings(t *testing.T) {
	as := assert.New(t)
	arr := ArrayOf("a", "b")
	val, err := arr.Value()
	as.NoError(err)
	as.Equal(`{"a","b"}`, val)
}

func TestArrayOf_LocalDate(t *testing.T) {
	as := assert.New(t)
	d := LocalDateOfYmd(2024, 6, 18)
	arr := ArrayOf(d)
	val, err := arr.Value()
	as.NoError(err)
	as.Equal("{2024-06-18}", val)
}

func TestArrayOf_Mixed(t *testing.T) {
	as := assert.New(t)
	arr := ArrayOf(1, "hello")
	val, err := arr.Value()
	as.NoError(err)
	as.Equal(`{1,"hello"}`, val)
}
