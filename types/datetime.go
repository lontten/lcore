package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgtype"
	"time"
)

// datetime
type DateTime struct {
	time.Time
}

func NowDateTime() DateTime {
	return DateTime{time.Now()}
}

func NowDateTimeP() *DateTime {
	return &DateTime{time.Now()}
}

func (t DateTime) ToString() string {
	return t.Format(`2006-01-02 15:04:05`)
}
func (t DateTime) ToDate() Date {
	return Date{t.ToGoTime()}
}
func (t DateTime) ToDateP() *Date {
	return &Date{t.ToGoTime()}
}
func (t DateTime) MarshalJSON() ([]byte, error) {
	tune := t.Format(`"2006-01-02 15:04:05"`)
	return []byte(tune), nil
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	*t = DateTime{Time: now}
	return err
}

// Value insert timestamp into mysql need this function.
func (t DateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof jstime.Time
func (t *DateTime) Scan(v any) error {
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
		*t = DateTime{v}
	case DateTime:
		*t = v
		return nil
	default:
		return fmt.Errorf("can not convert %v to DateTime", v)
	}
	if len(s) < 19 {
		return nil
	}
	now, err := time.ParseInLocation(`2006-01-02 15:04:05`, s, time.Local)
	if err != nil {
		return err
	}
	*t = DateTime{Time: now}
	return nil
}

func (t DateTime) ToGoTime() time.Time {
	return time.Unix(t.Unix(), 0)
}

type DateTimeList []DateTime

// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p DateTimeList) Value() (driver.Value, error) {
	var k []DateTime
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
func (p *DateTimeList) Scan(data any) error {
	array := pgtype.TimestampArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []DateTime
	list = make([]DateTime, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = DateTime{element.Time}
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}
func (t DateTime) Add(d *DurationOption) DateTime {
	if d == nil {
		return t
	}
	return DateTime{t.AddDate(d.year, d.month, d.day).Add(
		time.Duration(d.hour)*time.Hour +
			time.Duration(d.min)*time.Minute +
			time.Duration(d.sec)*time.Second +
			time.Duration(d.nsec)*time.Nanosecond,
	)}
}

// datetime
type AutoDateTime struct {
	time.Time
}

func (t AutoDateTime) MarshalJSON() ([]byte, error) {
	var tune string
	if t.Year() == 0 && t.Month() == time.January && t.Day() == 1 {
		tune = t.Format(`"15:04:05"`)
	} else {
		tune = t.Format(`"2006-01-02"`)
	}
	return []byte(tune), nil
}

func (t *AutoDateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	*t = AutoDateTime{Time: now}
	return err
}

// Value insert timestamp into mysql need this function.
func (t AutoDateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof jstime.Time
func (t *AutoDateTime) Scan(v any) error {
	var s = ""
	switch v := v.(type) {
	case string:
		s = v[:8]
	case []byte:
		s = string(v)[:8]
	case time.Time:
		*t = AutoDateTime{v}
	case Time:
		*t = AutoDateTime{v.Time}
	case Date:
		*t = AutoDateTime{v.Time}
	case AutoDateTime:
		*t = v
		return nil
	default:
		return fmt.Errorf("can not convert %v to types.AutoDateTime", v)
	}
	now, err := time.ParseInLocation(`2006-01-02 15:04:05`, s, time.Local)
	if err != nil {
		return err
	}
	*t = AutoDateTime{Time: now}
	return nil
}
