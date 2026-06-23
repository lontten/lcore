package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntList_ValueScan(t *testing.T) {
	as := assert.New(t)
	list := IntList{1, 2, -3}

	value, err := list.Value()
	as.NoError(err)

	var scanned IntList
	as.NoError(scanned.Scan(value))
	as.Equal(list, scanned)
}

func TestIntList_Empty(t *testing.T) {
	as := assert.New(t)
	empty := IntList{}
	val, err := empty.Value()
	as.NoError(err)
	as.Equal("{}", val)

	var scanned IntList
	as.NoError(scanned.Scan(val))
	as.Empty(scanned)
}

func TestIntList_ScanNil(t *testing.T) {
	as := assert.New(t)
	var list IntList
	list = IntList{1}
	as.NoError(list.Scan(nil))
	as.Equal(IntList{1}, list)
}
