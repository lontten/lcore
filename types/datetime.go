package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
)

// LocalDateTime 本地时区 日期时间
type LocalDateTime struct {
	data time.Time
}

// ----------------------- now ----------------------------

// NowDateTime 当前日期时间
func NowDateTime() LocalDateTime {
	return LocalDateTime{time.Now()}
}

func NowDateTimeP() *LocalDateTime {
	return &LocalDateTime{time.Now()}
}

// ----------------------- of ----------------------------

func LocalDateTimeOf(v time.Time) LocalDateTime {
	t := v.In(time.Local)
	return LocalDateTime{t}
}
func LocalDateTimePOf(v time.Time) *LocalDateTime {
	dateOnly := LocalDateTimeOf(v)
	return &dateOnly
}
func LocalDateTimeOfYmdHms(year, month, day, hour, min, sec int) LocalDateTime {
	dateOnly := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local)
	return LocalDateTime{dateOnly}
}
func LocalDateTimePOfYmdHms(year, month, day, hour, min, sec int) *LocalDateTime {
	dateOnly := LocalDateTimeOfYmdHms(year, month, day, hour, min, sec)
	return &dateOnly
}

// ----------------------- base ----------------------------

func (t LocalDateTime) IsZero() bool {
	return t.data.IsZero()
}

func (t LocalDateTime) String() string {
	return t.data.Format(`2006-01-02 15:04:05`)
}

func (t LocalDateTime) Add(d *DurationOption) LocalDateTime {
	if d == nil {
		return t
	}
	return LocalDateTime{t.data.AddDate(d.year, d.month, d.day).Add(
		time.Duration(d.hour)*time.Hour +
			time.Duration(d.min)*time.Minute +
			time.Duration(d.sec)*time.Second +
			time.Duration(d.nsec)*time.Nanosecond,
	)}
}

// ----------------------- comp ----------------------------

// Before
// t<d 返回true
// t>=d 返回false
func (t LocalDateTime) Before(d LocalDateTime) bool {
	return t.ToGoTime().Before(d.ToGoTime())
}

// After
// t>d 返回true
// t<=d 返回false
func (t LocalDateTime) After(d LocalDateTime) bool {
	return t.ToGoTime().After(d.ToGoTime())
}

// Eq 相等
// t==d 返回true
// t!=d 返回false
func (t LocalDateTime) Eq(d LocalDateTime) bool {
	return t.ToGoTime() == d.ToGoTime()
}

// ----------------------- to ----------------------------

func (t LocalDateTime) ToGoTime() time.Time {
	return t.data
}
func (t LocalDateTime) ToDate() LocalDate {
	return LocalDate{t.ToGoTime()}
}
func (t LocalDateTime) ToDateP() *LocalDate {
	return &LocalDate{t.ToGoTime()}
}

// ----------------------- parse ----------------------------

func LocalDateTimeParse(data string) (LocalDateTime, error) {
	now, err := time.ParseInLocation(`2006-01-02 15:04:05`, data, time.Local)
	return LocalDateTime{now}, err
}
func LocalDateTimeParseMust(data string) LocalDateTime {
	localTime, err := LocalDateTimeParse(data)
	if err != nil {
		panic(err)
	}
	return localTime
}

func LocalDateTimeParseMustP(data string) *LocalDateTime {
	localTime, err := LocalDateTimeParse(data)
	if err != nil {
		panic(err)
	}
	return &localTime
}

// ----------------------- json ----------------------------

func (t LocalDateTime) MarshalJSON() ([]byte, error) {
	tune := t.data.Format(`"2006-01-02 15:04:05"`)
	return []byte(tune), nil
}

func (t *LocalDateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	*t = LocalDateTime{data: now}
	return err
}

// ----------------------- db ----------------------------

// Value insert timestamp into mysql need this function.
func (t LocalDateTime) Value() (driver.Value, error) {
	return t.data, nil
}

// Scan valueof jstime.LocalTime
func (t *LocalDateTime) Scan(v any) error {
	if v == nil {
		return nil
	}
	var s = ""
	switch v := v.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	case time.Time:
		*t = LocalDateTimeOf(v)
	case LocalDateTime:
		*t = v
		return nil
	default:
		return fmt.Errorf("can not convert %v to LocalDateTime", v)
	}
	if len(s) < 19 {
		return nil
	}
	now, err := time.ParseInLocation(`2006-01-02 15:04:05`, s, time.Local)
	if err != nil {
		return err
	}
	*t = LocalDateTime{data: now}
	return nil
}

// ----------------------- list ----------------------------

type LocalDateTimeList []LocalDateTime

// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p LocalDateTimeList) Value() (driver.Value, error) {
	var k []LocalDateTime
	k = p
	marshal, err := json.Marshal(k)
	if err != nil {
		return nil, err
	}
	var s = string(marshal)
	if s != "null" {
		s = s[:0] + "{" + s[1:len(s)-1] + "}" + s[len(s):]
	} else {
		s = "{}"
	}
	return s, nil
}

// Scan 实现方法
func (p *LocalDateTimeList) Scan(data any) error {
	array := pgtype.TimestampArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []LocalDateTime
	list = make([]LocalDateTime, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = LocalDateTime{element.Time}
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}
