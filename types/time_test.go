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
	date := dateEnd.Time.AddDate(0, 0, 1)
	of := LocalDateOf(date)
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
	t.Log(zero)

	localTime := LocalTimeParseMust("00:00:00")
	as.Equal(zero, localTime)
}
