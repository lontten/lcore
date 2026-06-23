package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolList(t *testing.T) {
	as := assert.New(t)

	list := BoolList{true, false, true}

	value, err := list.Value()
	as.NoError(err)

	var scanned BoolList
	as.NoError(scanned.Scan(value))
	as.Equal(BoolList{true, false, true}, scanned)

	allFalse := BoolList{false, false}
	value2, err := allFalse.Value()
	as.NoError(err)
	var scannedFalse BoolList
	as.NoError(scannedFalse.Scan(value2))
	as.Equal(allFalse, scannedFalse)
}
