package types

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNowTime(t *testing.T) {
	now := time.Now()

	dateEnd := LocalDateOf(now)
	of := dateEnd.Add(Duration().Day(1))
	var dateTimeEnd = of.ToDateTimeP()

	fmt.Println(dateTimeEnd.String())
}

func TestLocalTime_String(t *testing.T) {
	as := assert.New(t)
	localTime := LocalTimeOfHms(10, 10, 10)
	s := localTime.String()
	as.Equal("10:10:10", s)
}

func TestLocalTimeParse(t *testing.T) {
	as := assert.New(t)
	localTime, err := LocalTimeParse("10:10:10")
	as.Nil(err)
	localTime2 := LocalTimeOfHms(10, 10, 10)
	as.Equal(localTime2, localTime)

	date := time.Date(1, time.January, 1, 10, 10, 10, 0, time.Local)
	localTime3 := LocalTimeOf(date)
	as.Equal(localTime3, localTime)
}

func TestLocalTimeZero(t *testing.T) {
	as := assert.New(t)
	zero := LocalTime{}
	as.True(zero.IsZero())
	as.Equal(zero.String(), "00:00:00")

	localTime := LocalTimeParseMust("00:00:00")
	as.NotEqual(zero, localTime)

	as.Equal(localTime, LocalTimeOfHms(0, 0, 0))
	as.Equal(localTime, LocalTimeOf(time.Date(0, time.January, 1, 0, 0, 0, 0, time.Local)))
	as.Equal(localTime.String(), "00:00:00")

	var n = LocalTime{time.Now()}
	as.False(n.IsZero())
	as.NotEqual(zero, n)
}

func TestLocalTimeList_Scan3(t *testing.T) {
	as := assert.New(t)
	var l = make([]LocalTime, 0)
	l = append(l, LocalTimeParseMust("00:00:00"))
	l = append(l, LocalTimeParseMust("00:00:01"))
	l = append(l, LocalTime{})

	var list LocalTimeList = l

	value, err := list.Value()
	as.NoError(err)

	var p LocalTimeList
	as.NoError(p.Scan(value))
	as.Equal(3, len(p))
	as.Equal(LocalTimeParseMust("00:00:00"), p[0])
	as.Equal(LocalTimeParseMust("00:00:01"), p[1])
	as.Equal(LocalTimeParseMust("00:00:00"), p[2])
	as.True(p[0].IsZero())
	as.True(p[2].IsZero())

}
