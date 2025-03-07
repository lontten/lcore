package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgtype"
	"time"
)

// date
type Date struct {
	time.Time
}

func NowDate() Date {
	return Date{time.Now()}
}

func NowDateP() *Date {
	return &Date{time.Now()}
}
func DateOf(t time.Time) Date {
	dateOnly := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return Date{dateOnly}
}
func DatePOf(t time.Time) *Date {
	dateOnly := DateOf(t)
	return &dateOnly
}
func DateOfYmd(year, month, day int) Date {
	dateOnly := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return Date{dateOnly}
}
func DatePOfYmd(year, month, day int) *Date {
	dateOnly := DateOfYmd(year, month, day)
	return &dateOnly
}
func (t Date) ToGoTime() time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
func (t Date) ToDateTime() DateTime {
	return DateTime{t.ToGoTime()}
}
func (t Date) ToDateTimeP() *DateTime {
	return &DateTime{t.ToGoTime()}
}

func (t Date) ToString() string {
	return t.Format(`2006-01-02`)
}

// Before
// t<d 返回true
// t>=d 返回false
func (t Date) Before(d Date) bool {
	return t.ToGoTime().Before(d.ToGoTime())
}

// After
// t>d 返回true
// t<=d 返回false
func (t Date) After(d Date) bool {
	return t.ToGoTime().After(d.ToGoTime())
}

func (t Date) Add(d *DurationOption) Date {
	if d == nil {
		return t
	}
	return Date{t.AddDate(d.year, d.month, d.day)}
}
func (d Date) AddTime(t Time) DateTime {
	return DateTime{time.Date(
		d.Time.Year(),
		d.Time.Month(),
		d.Time.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(), 0, nil,
	)}
}

func (t Date) MarshalJSON() ([]byte, error) {
	tune := t.Format(`"2006-01-02"`)
	return []byte(tune), nil
}

func (t *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"2006-01-02"`, string(data), time.Local)
	*t = Date{Time: now}
	return err
}

// Value insert timestamp into mysql need this function.
func (t Date) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Format("2006-01-02"), nil
}

// Scan valueof jstime.Time
func (t *Date) Scan(v any) error {
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
		*t = Date{v}
	case Date:
		*t = v
		return nil
	default:
		return fmt.Errorf("can not convert %v to Date", v)
	}
	if len(s) < 10 {
		return nil
	}
	now, err := time.ParseInLocation(`2006-01-02`, s[:10], time.Local)
	if err != nil {
		return err
	}
	*t = Date{Time: now}
	return nil
}

type DateList []Date

// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p DateList) Value() (driver.Value, error) {
	var k []Date
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
func (p *DateList) Scan(data any) error {
	array := pgtype.TimestampArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []Date
	list = make([]Date, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = Date{element.Time}
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}
