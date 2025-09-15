package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalDateTimeList_Scan(t *testing.T) {
	as := assert.New(t)
	var l = make([]LocalDateTime, 0)
	l = append(l, LocalDateTimeParseMust("0001-01-01 00:00:00"))
	l = append(l, LocalDateTimeParseMust("0001-01-01 00:00:01"))
	l = append(l, LocalDateTime{})

	var list LocalDateTimeList = l

	value, err := list.Value()
	as.NoError(err)

	var p LocalDateTimeList
	as.NoError(p.Scan(value))
	as.Equal(3, len(p))
	as.Equal(LocalDateTimeParseMust("0001-01-01 00:00:00"), p[0])
	as.Equal(LocalDateTimeParseMust("0001-01-01 00:00:01"), p[1])
	as.Equal(LocalDateTimeParseMust("0001-01-01 00:00:00"), p[2])
	as.True(p[0].IsZero())
	as.True(p[2].IsZero())
}
