package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
)

// LocalDate 本地时区日期
// 零值是 0001-01-01 00:00:00
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

func LocalDateOf(v time.Time) LocalDate {
	dateOnly := time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, time.Local)
	return LocalDate{dateOnly}
}
func LocalDatePOf(t time.Time) *LocalDate {
	dateOnly := LocalDateOf(t)
	return &dateOnly
}

// ----------------------- of Loc----------------------------

func LocalDateOfLoc(v time.Time) LocalDate {
	t := v.In(time.Local)
	dateOnly := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	return LocalDate{dateOnly}
}
func LocalDatePOfLoc(t time.Time) *LocalDate {
	dateOnly := LocalDateOfLoc(t)
	return &dateOnly
}

// ----------------------- of Ymd ----------------------------
func LocalDateOfYmd(year, month, day int) LocalDate {
	dateOnly := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return LocalDate{dateOnly}
}
func LocalDatePOfYmd(year, month, day int) *LocalDate {
	dateOnly := LocalDateOfYmd(year, month, day)
	return &dateOnly
}

// ----------------------- base ----------------------------
func (t LocalDate) String() string {
	return t.data.Format(`2006-01-02`)
}
func (t LocalDate) IsZero() bool {
	return t.String() == "0001-01-01"
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
	localTime, err := time.ParseInLocation(`2006-01-02`, data, time.Local)
	return LocalDate{localTime}, err
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
	localTime, err := time.ParseInLocation(`"2006-01-02"`, string(data), time.Local)
	*t = LocalDate{data: localTime}
	return err
}

// ----------------------- db ----------------------------

func (t LocalDate) Value() (driver.Value, error) {
	return t.data.Format("2006-01-02"), nil
}

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
		return nil
	case LocalDate:
		*t = v
		return nil
	default:
		return fmt.Errorf("can not convert %v to LocalDate", v)
	}
	if len(s) < 10 {
		return fmt.Errorf("can not convert %v to LocalDate", v)
	}
	localTime, err := time.ParseInLocation(`2006-01-02`, s[:10], time.Local)
	if err != nil {
		return err
	}
	*t = LocalDate{data: localTime}
	return nil
}

// ----------------------- list ----------------------------

type LocalDateList []LocalDate

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
	array := pgtype.DateArray{}
	if err := array.Scan(data); err != nil {
		return err
	}
	list := make([]LocalDate, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = LocalDateOf(element.Time)
	}
	*p = list
	return nil
}
