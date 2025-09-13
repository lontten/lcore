package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
)

// LocalTime 本地时间
// 零值是 0001-01-01 00:00:00
// 有效数值是 0000-01-01 xx:xx:xx
// "00:00:00" 转成 LocalTime 是 0000-01-01 00:00:00，和time零值 0001-01-01 00:00:00不同
// 只要确保 在 非默认值零值时，和 time零值 0001-01-01 00:00:00不同，即可
type LocalTime struct {
	data time.Time
}

// ----------------------- now ----------------------------

func NowLocalTime() LocalTime {
	return LocalTime{time.Now()}
}

func NowLocalTimeP() *LocalTime {
	return &LocalTime{time.Now()}
}

// ----------------------- of ----------------------------

func LocalTimeOf(v time.Time) LocalTime {
	t := v.In(time.Local)
	timeOnly := time.Date(0, time.January, 1, t.Hour(), t.Minute(), t.Second(), 0, time.Local)
	return LocalTime{timeOnly}
}
func LocalTimePOf(t time.Time) *LocalTime {
	timeOnly := LocalTimeOf(t)
	return &timeOnly
}
func LocalTimeOfHms(hour, min, sec int) LocalTime {
	timeOnly := time.Date(0, time.January, 1, hour, min, sec, 0, time.Local)
	return LocalTime{timeOnly}
}
func LocalTimePOfHms(hour, min, sec int) *LocalTime {
	timeOnly := LocalTimeOfHms(hour, min, sec)
	return &timeOnly
}

// ----------------------- base ----------------------------

func (t LocalTime) IsZero() bool {
	return t.data.IsZero()
}

func (t LocalTime) String() string {
	return t.data.Format(`15:04:05`)
}

func (t LocalTime) Add(d *DurationOption) LocalTime {
	if d == nil {
		return t
	}
	return LocalTime{t.ToGoTime().Add(
		time.Duration(d.hour)*time.Hour +
			time.Duration(d.min)*time.Minute +
			time.Duration(d.sec)*time.Second +
			time.Duration(d.nsec)*time.Nanosecond,
	)}
}
func (t LocalTime) AddData(d LocalDate) LocalDateTime {
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
func (t LocalTime) Before(d LocalTime) bool {
	return t.ToGoTime().Before(d.ToGoTime())
}

// After
// t>d 返回true
// t<=d 返回false
func (t LocalTime) After(d LocalTime) bool {
	return t.ToGoTime().After(d.ToGoTime())
}

// Eq
// t==d 返回true
// t!=d 返回false
func (t LocalTime) Eq(d LocalTime) bool {
	return t.ToGoTime() == d.ToGoTime()
}

// ----------------------- to ----------------------------
func (t LocalTime) ToGoTime() time.Time {
	return time.Date(0, time.January, 1, t.data.Hour(), t.data.Minute(), t.data.Second(), 0, time.Local)
}
func (t LocalTime) ToLocalDateTime() LocalDateTime {
	return LocalDateTime{t.ToGoTime()}
}
func (t LocalTime) ToLocalDateTimeP() *LocalDateTime {
	return &LocalDateTime{t.ToGoTime()}
}

// ----------------------- parse ----------------------------

// LocalTimeParse 解析结果为 0年1月1日 00:00:00，和time零值的 1年1月1日 00:00:00不同
func LocalTimeParse(data string) (LocalTime, error) {
	now, err := time.ParseInLocation(`15:04:05`, data, time.Local)
	return LocalTime{now}, err
}
func LocalTimeParseMust(data string) LocalTime {
	localTime, err := LocalTimeParse(data)
	if err != nil {
		panic(err)
	}
	return localTime
}

func LocalTimeParseMustP(data string) *LocalTime {
	localTime, err := LocalTimeParse(data)
	if err != nil {
		panic(err)
	}
	return &localTime
}

// ----------------------- json ----------------------------

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tune := t.data.Format(`"15:04:05"`)
	return []byte(tune), nil
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"15:04:05"`, string(data), time.Local)
	*t = LocalTime{now}
	return err
}

// ----------------------- db ----------------------------

func (t LocalTime) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.data.Format("15:04:05"), nil
}

func (t *LocalTime) Scan(v any) error {
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
		*t = LocalTimeOf(v)
	case LocalTime:
		*t = v
		return nil
	default:
		return fmt.Errorf("can not convert %v to LocalTime", v)
	}
	if len(s) < 8 {
		return nil
	}
	now, err := time.ParseInLocation(`15:04:05`, s[len(s)-8:], time.Local)
	if err != nil {
		return err
	}
	*t = LocalTime{data: now}
	return nil
}

// ----------------------- list ----------------------------

type LocalTimeList []LocalTime

// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p LocalTimeList) Value() (driver.Value, error) {
	var k []LocalTime
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
func (p *LocalTimeList) Scan(data any) error {
	array := pgtype.TimestampArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []LocalTime
	list = make([]LocalTime, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = LocalTimeOf(element.Time)
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}
