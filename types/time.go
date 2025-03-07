package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgtype"
	"time"
)

type Time struct {
	time.Time
}

func NowTime() Time {
	return Time{time.Now()}
}

func NowTimeP() *Time {
	return &Time{time.Now()}
}
func TimeOf(t time.Time) Time {
	timeOnly := time.Date(0, 0, 0, t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	return Time{timeOnly}
}
func TimePOf(t time.Time) *Time {
	timeOnly := TimeOf(t)
	return &timeOnly
}
func TimeOfHms(hour, min, sec int) Time {
	timeOnly := time.Date(0, 0, 0, hour, min, sec, 0, time.Local)
	return Time{timeOnly}
}
func TimePOfHms(hour, min, sec int) *Time {
	timeOnly := TimeOfHms(hour, min, sec)
	return &timeOnly
}
func (t Time) ToGoTime() time.Time {
	return time.Date(0, 0, 0, t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}
func (t Time) ToDateTime() DateTime {
	return DateTime{t.ToGoTime()}
}
func (t Time) ToDateTimeP() *DateTime {
	return &DateTime{t.ToGoTime()}
}
func (t Time) ToString() string {
	return t.Time.Format(`15:04:05`)
}
func (t Time) Add(d *DurationOption) Time {
	if d == nil {
		return t
	}
	return Time{t.ToGoTime().Add(
		time.Duration(d.hour)*time.Hour +
			time.Duration(d.min)*time.Minute +
			time.Duration(d.sec)*time.Second +
			time.Duration(d.nsec)*time.Nanosecond,
	)}
}
func (t Time) AddData(d Date) DateTime {
	return DateTime{time.Date(
		d.Time.Year(),
		d.Time.Month(),
		d.Time.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(), 0, time.Local,
	)}
}

func (t Time) MarshalJSON() ([]byte, error) {
	tune := t.Time.Format(`"15:04:05"`)
	return []byte(tune), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"15:04:05"`, string(data), time.Local)
	*t = Time{now}
	return err
}

func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Format("15:04:05"), nil
}

func (t *Time) Scan(v any) error {
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
		*t = Time{v}
	case Time:
		*t = v
		return nil
	default:
		return fmt.Errorf("can not convert %v to Time", v)
	}
	if len(s) < 8 {
		return nil
	}
	now, err := time.ParseInLocation(`15:04:05`, s[len(s)-8:], time.Local)
	if err != nil {
		return err
	}
	*t = Time{Time: now}
	return nil
}

type TimeList []Time

// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p TimeList) Value() (driver.Value, error) {
	var k []Time
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
func (p *TimeList) Scan(data any) error {
	array := pgtype.TimestampArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []Time
	list = make([]Time, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = Time{element.Time}
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}
