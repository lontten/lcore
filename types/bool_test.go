package types

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolList_ValueScan(t *testing.T) {
	as := assert.New(t)
	list := BoolList{true, false, true}

	value, err := list.Value()
	as.NoError(err)

	var scanned BoolList
	as.NoError(scanned.Scan(value))
	as.Equal(BoolList{true, false, true}, scanned)
}

func TestBoolList_AllFalse(t *testing.T) {
	as := assert.New(t)
	allFalse := BoolList{false, false}
	value, err := allFalse.Value()
	as.NoError(err)

	var scanned BoolList
	as.NoError(scanned.Scan(value))
	as.Equal(allFalse, scanned)
}

func TestBoolList_EmptyAndNil(t *testing.T) {
	as := assert.New(t)
	empty := BoolList{}
	val, err := empty.Value()
	as.NoError(err)
	as.Equal("{}", val)

	var scanned BoolList
	as.NoError(scanned.Scan(nil))
	as.Nil(scanned)
}

func TestBoolList_Sort(t *testing.T) {
	as := assert.New(t)
	list := BoolList{true, false, true, false}
	sort.Sort(list)
	as.Equal(BoolList{false, false, true, true}, list)
}

func TestBoolList_ScanPgLiteral(t *testing.T) {
	as := assert.New(t)
	var list BoolList
	as.NoError(list.Scan("{true,false}"))
	as.Equal(BoolList{true, false}, list)
}
