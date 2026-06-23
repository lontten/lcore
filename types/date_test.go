package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalDate_Zero(t *testing.T) {
	as := assert.New(t)

	zero := LocalDate{}
	as.True(zero.IsZero())
	as.Equal("0001-01-01", zero.String())

	fromParse := LocalDateParseMust("0001-01-01")
	as.True(fromParse.IsZero())
	as.Equal("2024-06-18", LocalDateOfYmd(2024, 6, 18).String())
	as.False(LocalDateOfYmd(2024, 6, 18).IsZero())
}

func TestLocalDate_Constructors(t *testing.T) {
	as := assert.New(t)

	d := LocalDateOfYmd(2024, 6, 18)
	as.Equal("2024-06-18", d.String())

	tm := time.Date(2024, 6, 18, 15, 30, 0, 0, time.Local)
	as.Equal(d, LocalDateOf(tm))

	utc := time.Date(2024, 6, 18, 0, 0, 0, 0, time.UTC)
	as.Equal(LocalDateOfYmd(2024, 6, 18), LocalDateOfLoc(utc))

	p := LocalDatePOfYmd(2024, 1, 1)
	require.NotNil(t, p)
	as.Equal("2024-01-01", p.String())

	nowP := NowDateP()
	require.NotNil(t, nowP)
	as.Equal(NowDate().String(), nowP.String())
}

func TestLocalDate_Parse(t *testing.T) {
	as := assert.New(t)

	got, err := LocalDateParse("2023-01-01")
	as.NoError(err)
	as.Equal(LocalDateParseMust("2023-01-01"), got)

	_, err = LocalDateParse("invalid")
	as.Error(err)
}

func TestLocalDate_ParseMust(t *testing.T) {
	as := assert.New(t)
	as.Panics(func() { LocalDateParseMust("bad") })
	as.Panics(func() { LocalDateParseMustP("bad") })
	p := LocalDateParseMustP("2023-01-01")
	require.NotNil(t, p)
	as.Equal("2023-01-01", p.String())
}

func TestLocalDate_Compare(t *testing.T) {
	as := assert.New(t)
	early := LocalDateOfYmd(2020, 1, 1)
	late := LocalDateOfYmd(2021, 1, 1)
	as.True(early.Before(late))
	as.True(late.After(early))
	as.True(early.Eq(LocalDateOfYmd(2020, 1, 1)))
}

func TestLocalDate_Add(t *testing.T) {
	as := assert.New(t)
	d := LocalDateOfYmd(2024, 1, 15)
	as.Equal("2024-02-15", d.Add(Duration().Month(1)).String())
	as.Equal(d, d.Add(nil))
	as.Equal(
		"2024-06-18 10:30:00",
		LocalDateOfYmd(2024, 6, 18).AddTime(LocalTimeOfHms(10, 30, 0)).String(),
	)
}

func TestLocalDate_JSON(t *testing.T) {
	as := assert.New(t)
	src := LocalDateOfYmd(2024, 6, 18)
	data, err := src.MarshalJSON()
	as.NoError(err)
	as.Equal(`"2024-06-18"`, string(data))

	var dst LocalDate
	as.NoError(json.Unmarshal(data, &dst))
	as.Equal(src.String(), dst.String())

	existing := LocalDateOfYmd(2020, 1, 1)
	as.NoError(json.Unmarshal([]byte("null"), &existing))
	as.Equal("2020-01-01", existing.String())
}

func TestLocalDate_DB(t *testing.T) {
	as := assert.New(t)
	src := LocalDateOfYmd(2024, 6, 18)
	val, err := src.Value()
	as.NoError(err)
	as.Equal("2024-06-18", val)

	var got LocalDate
	as.NoError(got.Scan("2024-06-18"))
	as.Equal(src, got)

	as.NoError(got.Scan([]byte("2024-06-18")))
	as.NoError(got.Scan(time.Date(2024, 6, 18, 0, 0, 0, 0, time.Local)))
	as.NoError(got.Scan(src))

	existing := LocalDateOfYmd(2020, 1, 1)
	as.NoError(existing.Scan(nil))
	as.Equal("2020-01-01", existing.String())

	as.Error(got.Scan(123))
	as.Error(got.Scan("short"))
}

func TestLocalDateList_ValueScan(t *testing.T) {
	as := assert.New(t)
	list := LocalDateList{
		LocalDateParseMust("2023-01-01"),
		LocalDate{},
		LocalDateParseMust("1001-01-01"),
	}

	value, err := list.Value()
	as.NoError(err)

	var scanned LocalDateList
	as.NoError(scanned.Scan(value))
	as.Equal(3, len(scanned))
	as.Equal(LocalDateParseMust("2023-01-01"), scanned[0])
	as.True(scanned[1].IsZero())
	as.Equal(LocalDateParseMust("1001-01-01"), scanned[2])
}

func TestLocalDate_AddTimeToZero(t *testing.T) {
	as := assert.New(t)
	hms := LocalTimeOfHms(0, 0, 0).ToLocalDateTime().ToDate()
	as.True(hms.IsZero())
}
