package types

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringList_ValueScan(t *testing.T) {
	as := assert.New(t)
	list := StringList{"hello", "world", "a,b"}

	value, err := list.Value()
	as.NoError(err)

	var scanned StringList
	as.NoError(scanned.Scan(value))
	as.Equal(list, scanned)
}

func TestStringList_SpecialChars(t *testing.T) {
	as := assert.New(t)
	list := StringList{`quote"here`, `back\slash`}

	value, err := list.Value()
	as.NoError(err)

	var scanned StringList
	as.NoError(scanned.Scan(value))
	as.Equal(list, scanned)
}

func TestStringList_EmptyAndNil(t *testing.T) {
	as := assert.New(t)
	empty := StringList{}
	val, err := empty.Value()
	as.NoError(err)
	as.Equal("{}", val)

	var scanned StringList
	as.NoError(scanned.Scan(nil))
	as.Nil(scanned)
}

func TestStringList_Sort(t *testing.T) {
	as := assert.New(t)
	list := StringList{"c", "a", "b"}
	sort.Sort(list)
	as.Equal(StringList{"a", "b", "c"}, list)
}
