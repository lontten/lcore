package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalDateList_Scan3(t *testing.T) {
	as := assert.New(t)
	var l = make([]LocalDate, 0)
	l = append(l, LocalDateParseMust("2023-01-01"))
	l = append(l, LocalDate{})

	var list LocalDateList = l

	value, err := list.Value()
	as.NoError(err)

	var p LocalDateList
	as.NoError(p.Scan(value))
	as.Equal(LocalDateParseMust("2023-01-01"), p[0])
	as.Equal(LocalDateParseMust("0001-01-01"), p[1])
	as.True(p[1].IsZero())

	hms := LocalTimeOfHms(0, 0, 0).ToLocalDateTime().ToDate()
	as.True(hms.IsZero())
}

func TestLocalDateList_Scan4(t *testing.T) {
	as := assert.New(t)

	v1 := LocalTimeOfHms(0, 0, 0)
	as.Equal("00:00:00", v1.String())
	v2 := v1.ToLocalDateTime()
	as.Equal("0001-01-01 00:00:00", v2.String())
	v3 := v2.ToDate()
	as.Equal("0001-01-01", v3.String())

	as.True(v3.IsZero())
}

func TestLocalDateList_Scan(t *testing.T) {
	as := assert.New(t)
	var l = make([]LocalDate, 0)
	l = append(l, LocalDateParseMust("0001-01-01"))
	l = append(l, LocalDateParseMust("1001-01-01"))
	l = append(l, LocalDate{})

	var list LocalDateList = l

	value, err := list.Value()
	as.NoError(err)

	var p LocalDateList
	as.NoError(p.Scan(value))
	as.Equal(3, len(p))
	as.Equal(LocalDateParseMust("0001-01-01"), p[0])
	as.Equal(LocalDateParseMust("1001-01-01"), p[1])
	as.Equal(LocalDateParseMust("0001-01-01"), p[2])
	as.True(p[0].IsZero())
	as.True(p[2].IsZero())
}
