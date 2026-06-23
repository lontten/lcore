package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalDateTime_Zero(t *testing.T) {
	as := assert.New(t)
	zero := LocalDateTime{}
	as.True(zero.IsZero())
	as.Equal("0001-01-01 00:00:00", zero.String())
}

func TestLocalDateTime_Constructors(t *testing.T) {
	as := assert.New(t)
	dt := LocalDateTimeOfYmdHms(2024, 6, 18, 15, 4, 5)
	as.Equal("2024-06-18 15:04:05", dt.String())

	tm := time.Date(2024, 6, 18, 15, 4, 5, 0, time.Local)
	as.Equal(dt, LocalDateTimeOf(tm))

	utc := time.Date(2024, 6, 18, 12, 0, 0, 0, time.UTC)
	as.Equal(LocalDateTimeOf(utc.In(time.Local)), LocalDateTimeOfLoc(utc))

	p := LocalDateTimePOfYmdHms(2024, 1, 1, 0, 0, 0)
	require.NotNil(t, p)
	as.Equal("2024-01-01 00:00:00", p.String())
}

func TestLocalDateTime_Parse(t *testing.T) {
	as := assert.New(t)
	got, err := LocalDateTimeParse("2024-06-18 15:04:05")
	as.NoError(err)
	as.Equal(LocalDateTimeOfYmdHms(2024, 6, 18, 15, 4, 5), got)

	_, err = LocalDateTimeParse("bad")
	as.Error(err)
}

func TestLocalDateTime_ParseMust(t *testing.T) {
	as := assert.New(t)
	as.Panics(func() { LocalDateTimeParseMust("bad") })
	p := LocalDateTimeParseMustP("2024-06-18 15:04:05")
	require.NotNil(t, p)
	as.Equal("2024-06-18 15:04:05", p.String())
}

func TestLocalDateTime_Compare(t *testing.T) {
	as := assert.New(t)
	early := LocalDateTimeOfYmdHms(2024, 1, 1, 10, 0, 0)
	late := LocalDateTimeOfYmdHms(2024, 1, 1, 12, 0, 0)
	as.True(early.Before(late))
	as.True(late.After(early))
	as.True(early.Eq(LocalDateTimeOfYmdHms(2024, 1, 1, 10, 0, 0)))
}

func TestLocalDateTime_Add(t *testing.T) {
	as := assert.New(t)
	dt := LocalDateTimeOfYmdHms(2024, 1, 15, 10, 0, 0)
	as.Equal("2024-01-16 11:30:00", dt.Add(Duration().Day(1).Hour(1).Min(30)).String())
	as.Equal(dt, dt.Add(nil))
}

func TestLocalDateTime_Convert(t *testing.T) {
	as := assert.New(t)
	dt := LocalDateTimeOfYmdHms(2024, 6, 18, 15, 4, 5)
	as.Equal("2024-06-18", dt.ToDate().String())
	p := dt.ToDateP()
	require.NotNil(t, p)
	as.Equal("2024-06-18", p.String())
}

func TestLocalDateTime_JSON(t *testing.T) {
	as := assert.New(t)
	src := LocalDateTimeOfYmdHms(2024, 6, 18, 15, 4, 5)
	data, err := src.MarshalJSON()
	as.NoError(err)
	as.Equal(`"2024-06-18 15:04:05"`, string(data))

	var dst LocalDateTime
	as.NoError(json.Unmarshal(data, &dst))
	as.Equal(src, dst)

	existing := LocalDateTimeOfYmdHms(2020, 1, 1, 0, 0, 0)
	as.NoError(json.Unmarshal([]byte("null"), &existing))
	as.Equal("2020-01-01 00:00:00", existing.String())
}

func TestLocalDateTime_DB(t *testing.T) {
	as := assert.New(t)
	src := LocalDateTimeOfYmdHms(2024, 6, 18, 15, 4, 5)
	val, err := src.Value()
	as.NoError(err)
	as.Equal(src.ToGoTime(), val)

	var got LocalDateTime
	as.NoError(got.Scan("2024-06-18 15:04:05"))
	as.Equal(src, got)

	as.NoError(got.Scan([]byte("2024-06-18 15:04:05")))
	as.NoError(got.Scan(src.ToGoTime()))
	as.NoError(got.Scan(src))

	existing := src
	as.NoError(existing.Scan(nil))
	as.Equal(src, existing)

	as.Error(got.Scan(123))
	as.Error(got.Scan("short"))
}

func TestLocalDateTimeList_ValueScan(t *testing.T) {
	as := assert.New(t)
	list := LocalDateTimeList{
		LocalDateTimeParseMust("0001-01-01 00:00:00"),
		LocalDateTimeParseMust("0001-01-01 00:00:01"),
		LocalDateTime{},
	}

	value, err := list.Value()
	as.NoError(err)

	var scanned LocalDateTimeList
	as.NoError(scanned.Scan(value))
	as.Equal(3, len(scanned))
	as.True(scanned[0].IsZero())
	as.Equal(LocalDateTimeParseMust("0001-01-01 00:00:01"), scanned[1])
	as.True(scanned[2].IsZero())
}
