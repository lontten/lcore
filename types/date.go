package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
)

// LocalDate 本地时区日期
type LocalDate struct {
	data time.Time
}

// ----------------------- now ----------------------------

// NowLocalDate 当前时间
func NowLocalDate() LocalDate {
	return LocalDate{time.Now()}
}

func NowLocalDateP() *LocalDate {
	return &LocalDate{time.Now()}
}

// ----------------------- of ----------------------------

// LocalDateOf
// v 可能是 任意时区，需要先转成 本地时区，然后再转成 LocalDate
func LocalDateOf(v time.Time) LocalDate {
	t := v.In(time.Local)
	dateOnly := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	return LocalDate{dateOnly}
}
func LocalDatePOf(t time.Time) *LocalDate {
	dateOnly := LocalDateOf(t)
	return &dateOnly
}
func LocalDateOfYmd(year, month, day int) LocalDate {
	dateOnly := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return LocalDate{dateOnly}
}
func LocalDatePOfYmd(year, month, day int) *LocalDate {
	dateOnly := LocalDateOfYmd(year, month, day)
	return &dateOnly
}

// ----------------------- base ----------------------------

func (t LocalDate) IsZero() bool {
	return t.data.IsZero()
}

func (t LocalDate) String() string {
	return t.data.Format(`2006-01-02`)
}

func (t LocalDate) Add(d *DurationOption) LocalDate {
	if d == nil {
		return t
	}
	return LocalDate{t.data.AddDate(d.year, d.month, d.day)}
}
func (d LocalDate) AddTime(t LocalTime) LocalDateTime {
	return LocalDateTime{time.Date(
		d.data.Year(),
		d.data.Month(),
		d.data.Day(),
		t.data.Hour(),
		t.data.Minute(),
		t.data.Second(),
		0, time.Local,
	)}
}

// ----------------------- comp ----------------------------

// Before
// t<d 返回true
// t>=d 返回false
func (t LocalDate) Before(d LocalDate) bool {
	return t.ToGoTime().Before(d.ToGoTime())
}

// After
// t>d 返回true
// t<=d 返回false
func (t LocalDate) After(d LocalDate) bool {
	return t.ToGoTime().After(d.ToGoTime())
}

// Eq
// t==d 返回true
// t!=d 返回false
func (t LocalDate) Eq(d LocalDate) bool {
	return t.ToGoTime() == d.ToGoTime()
}

// ----------------------- to ----------------------------

func (t LocalDate) ToGoTime() time.Time {
	return time.Date(t.data.Year(), t.data.Month(), t.data.Day(), 0, 0, 0, 0, time.Local)
}
func (t LocalDate) ToDateTime() LocalDateTime {
	return LocalDateTime{t.ToGoTime()}
}
func (t LocalDate) ToDateTimeP() *LocalDateTime {
	return &LocalDateTime{t.ToGoTime()}
}

// ----------------------- parse ----------------------------

func LocalDateParse(data string) (LocalDate, error) {
	now, err := time.ParseInLocation(`2006-01-02`, data, time.Local)
	return LocalDate{now}, err
}
func LocalDateParseMust(data string) LocalDate {
	localTime, err := LocalDateParse(data)
	if err != nil {
		panic(err)
	}
	return localTime
}

func LocalDateParseMustP(data string) *LocalDate {
	localTime, err := LocalDateParse(data)
	if err != nil {
		panic(err)
	}
	return &localTime
}

// ----------------------- json ----------------------------

func (t LocalDate) MarshalJSON() ([]byte, error) {
	tune := t.data.Format(`"2006-01-02"`)
	return []byte(tune), nil
}

func (t *LocalDate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"2006-01-02"`, string(data), time.Local)
	*t = LocalDate{data: now}
	return err
}

// ----------------------- db ----------------------------

// Value insert timestamp into mysql need this function.
func (t LocalDate) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.data.Format("2006-01-02"), nil
}

// Scan valueof jstime.LocalTime
func (t *LocalDate) Scan(v any) error {
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
		*t = LocalDateOf(v)
	case LocalDate:
		*t = v
		return nil
	default:
		return fmt.Errorf("can not convert %v to LocalDate", v)
	}
	if len(s) < 10 {
		return nil
	}
	now, err := time.ParseInLocation(`2006-01-02`, s[:10], time.Local)
	if err != nil {
		return err
	}
	*t = LocalDate{data: now}
	return nil
}

// ----------------------- list ----------------------------

type LocalDateList []LocalDate

// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p LocalDateList) Value() (driver.Value, error) {
	var k []LocalDate
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
func (p *LocalDateList) Scan(data any) error {
	array := pgtype.TimestampArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []LocalDate
	list = make([]LocalDate, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = LocalDateOf(element.Time)
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}
