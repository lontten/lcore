package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDurationOption_Chain(t *testing.T) {
	as := assert.New(t)
	d := Duration().Year(1).Month(2).Day(3).Hour(4).Min(5).Sec(6).Nsec(7)
	as.NotNil(d)
}

func TestDurationOption_WithLocalDate(t *testing.T) {
	as := assert.New(t)
	d := LocalDateOfYmd(2024, 1, 15)
	as.Equal("2024-02-20", d.Add(Duration().Month(1).Day(5)).String())
}

func TestDurationOption_WithLocalTime(t *testing.T) {
	as := assert.New(t)
	tm := LocalTimeOfHms(10, 0, 0)
	as.Equal("11:15:30", tm.Add(Duration().Hour(1).Min(15).Sec(30)).String())
}

func TestDurationOption_WithLocalDateTime(t *testing.T) {
	as := assert.New(t)
	dt := LocalDateTimeOfYmdHms(2024, 1, 1, 10, 0, 0)
	as.Equal("2024-01-02 11:00:00", dt.Add(Duration().Day(1).Hour(1)).String())
}
