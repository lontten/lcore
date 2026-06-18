package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalTime_Zero(t *testing.T) {
	as := assert.New(t)

	zero := LocalTime{}
	as.True(zero.IsZero())
	as.Equal("00:00:00", zero.String())

	fromHms := LocalTimeOfHms(0, 0, 0)
	fromParse := LocalTimeParseMust("00:00:00")

	as.Equal(zero, fromHms)
	as.Equal(zero, fromParse)
	as.True(fromHms.IsZero())
	as.True(fromParse.IsZero())

	nonZero := LocalTimeOfHms(10, 0, 0)
	as.False(nonZero.IsZero())
}

func TestLocalTime_Constructors(t *testing.T) {
	as := assert.New(t)

	as.Equal("10:10:10", LocalTimeOfHms(10, 10, 10).String())

	date := time.Date(1, time.January, 1, 10, 10, 10, 0, time.Local)
	as.Equal(LocalTimeOfHms(10, 10, 10), LocalTimeOf(date))

	utc := time.Date(2024, 6, 18, 12, 0, 0, 0, time.UTC)
	as.Equal(LocalTimeOfHms(12, 0, 0), LocalTimeOf(utc))
	as.Equal(localTimeFrom(utc), LocalTimeOfLoc(utc))

	pOf := LocalTimePOf(date)
	require.NotNil(t, pOf)
	as.Equal(LocalTimeOfHms(10, 10, 10), *pOf)

	pOfHms := LocalTimePOfHms(10, 10, 10)
	require.NotNil(t, pOfHms)
	as.Equal(LocalTimeOfHms(10, 10, 10), *pOfHms)

	pOfLoc := LocalTimePOfLoc(utc)
	require.NotNil(t, pOfLoc)
	as.Equal(LocalTimeOfLoc(utc), *pOfLoc)

	nowP := NowTimeP()
	require.NotNil(t, nowP)
	as.Equal(NowTime(), *nowP)

	now := NowTime()
	expected := localTimeFrom(time.Now())
	as.Equal(expected.hour, now.hour)
	as.Equal(expected.minute, now.minute)
	as.Equal(expected.second, now.second)

	as.Panics(func() { LocalTimeOfHms(25, 0, 0) })
}

func TestLocalTime_Parse(t *testing.T) {
	as := assert.New(t)

	successCases := []struct {
		name  string
		input string
		want  LocalTime
	}{
		{"normal", "10:10:10", LocalTimeOfHms(10, 10, 10)},
		{"zero", "00:00:00", LocalTime{}},
		{"padding", "01:02:03", LocalTimeOfHms(1, 2, 3)},
	}
	for _, tc := range successCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := LocalTimeParse(tc.input)
			as.NoError(err)
			as.Equal(tc.want, got)

			gotP, err := LocalTimeParseP(tc.input)
			as.NoError(err)
			require.NotNil(t, gotP)
			as.Equal(tc.want, *gotP)
		})
	}

	failCases := []string{"", "abc", "25:00:00", "12:60:00", "10:10"}
	for _, input := range failCases {
		t.Run("fail_"+input, func(t *testing.T) {
			_, err := LocalTimeParse(input)
			as.Error(err)

			gotP, err := LocalTimeParseP(input)
			as.Error(err)
			as.Nil(gotP)
		})
	}
}

func TestLocalTime_ParseMust(t *testing.T) {
	as := assert.New(t)

	as.Panics(func() { LocalTimeParseMust("invalid") })
	as.Panics(func() { LocalTimeParseMustP("invalid") })

	p := LocalTimeParseMustP("10:10:10")
	require.NotNil(t, p)
	as.Equal(LocalTimeOfHms(10, 10, 10), *p)
}

func TestLocalTime_String(t *testing.T) {
	as := assert.New(t)
	as.Equal("01:02:03", LocalTimeOfHms(1, 2, 3).String())
}

func TestLocalTime_Compare(t *testing.T) {
	as := assert.New(t)

	early := LocalTimeOfHms(10, 0, 0)
	late := LocalTimeOfHms(12, 30, 45)
	same := LocalTimeOfHms(10, 0, 0)

	as.True(early.Before(late))
	as.False(late.Before(early))
	as.False(early.Before(same))

	as.True(late.After(early))
	as.False(early.After(late))
	as.False(early.After(same))

	as.True(early.Eq(same))
	as.False(early.Eq(late))
	as.True(early == same)
}

func TestLocalTime_Add(t *testing.T) {
	as := assert.New(t)

	as.Equal("12:00:00", LocalTimeOfHms(10, 0, 0).Add(Duration().Hour(2)).String())
	as.Equal("00:30:00", LocalTimeOfHms(23, 30, 0).Add(Duration().Hour(1)).String())

	original := LocalTimeOfHms(10, 0, 0)
	as.Equal(original, original.Add(nil))
}

func TestLocalTime_Convert(t *testing.T) {
	as := assert.New(t)

	tm := LocalTimeOfHms(15, 4, 5).ToGoTime()
	as.Equal(localTimeAnchorYear, tm.Year())
	as.Equal(time.January, tm.Month())
	as.Equal(1, tm.Day())
	as.Equal(15, tm.Hour())
	as.Equal(4, tm.Minute())
	as.Equal(5, tm.Second())

	as.Equal("0001-01-01 15:04:05", LocalTimeOfHms(15, 4, 5).ToLocalDateTime().String())

	dtP := LocalTimeOfHms(15, 4, 5).ToLocalDateTimeP()
	require.NotNil(t, dtP)
	as.Equal("0001-01-01 15:04:05", dtP.String())

	as.Equal(
		"2024-06-18 15:04:05",
		LocalDateOfYmd(2024, 6, 18).AddTime(LocalTimeOfHms(15, 4, 5)).String(),
	)
}

func TestLocalTime_JSON(t *testing.T) {
	as := assert.New(t)

	src := LocalTimeOfHms(10, 10, 10)
	data, err := src.MarshalJSON()
	as.NoError(err)
	as.Equal(`"10:10:10"`, string(data))

	var dst LocalTime
	as.NoError(json.Unmarshal(data, &dst))
	as.Equal(src, dst)

	existing := LocalTimeOfHms(12, 0, 0)
	as.NoError(json.Unmarshal([]byte("null"), &existing))
	as.Equal(LocalTimeOfHms(12, 0, 0), existing)
}

func TestLocalTime_DB(t *testing.T) {
	as := assert.New(t)

	src := LocalTimeOfHms(10, 10, 10)
	val, err := src.Value()
	as.NoError(err)
	as.Equal("10:10:10", val)

	cases := []struct {
		name  string
		input any
		want  LocalTime
	}{
		{"string", "10:10:10", LocalTimeOfHms(10, 10, 10)},
		{"bytes", []byte("10:10:10"), LocalTimeOfHms(10, 10, 10)},
		{"datetime_suffix", "2006-01-02 15:04:05", LocalTimeOfHms(15, 4, 5)},
		{"time", time.Date(1, time.January, 1, 8, 30, 0, 0, time.Local), LocalTimeOfHms(8, 30, 0)},
		{"local_time", LocalTimeOfHms(9, 9, 9), LocalTimeOfHms(9, 9, 9)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got LocalTime
			as.NoError(got.Scan(tc.input))
			as.Equal(tc.want, got)
		})
	}

	var invalid LocalTime
	as.Error(invalid.Scan("bad"))
	as.Error(invalid.Scan(123))

	existing := LocalTimeOfHms(12, 0, 0)
	as.NoError(existing.Scan(nil))
	as.Equal(LocalTimeOfHms(12, 0, 0), existing)
}

func TestLocalTimeList(t *testing.T) {
	as := assert.New(t)

	list := LocalTimeList{
		LocalTimeParseMust("00:00:00"),
		LocalTimeParseMust("00:00:01"),
		LocalTime{},
	}

	value, err := list.Value()
	as.NoError(err)

	var scanned LocalTimeList
	as.NoError(scanned.Scan(value))
	as.Equal(3, len(scanned))
	as.Equal(LocalTimeParseMust("00:00:00"), scanned[0])
	as.Equal(LocalTimeParseMust("00:00:01"), scanned[1])
	as.Equal(LocalTime{}, scanned[2])
	as.True(scanned[0].IsZero())
	as.True(scanned[2].IsZero())
}
