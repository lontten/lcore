package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// LocalDate 本地日历日期（仅年月日，不含时分秒与时区偏移）。
// 语义等同 java.time.LocalDate，统一使用 time.Local。
// 零值为 0001-01-01；IsZero 按字符串 "0001-01-01" 判断。
type LocalDate struct {
	data time.Time
}

// ----------------------- now ----------------------------

// NowDate 返回当前本地日期。
func NowDate() LocalDate {
	return LocalDate{time.Now()}
}

// NowDateP 返回当前本地日期的指针。
func NowDateP() *LocalDate {
	return &LocalDate{time.Now()}
}

// ----------------------- of ----------------------------

// LocalDateOf 从 time.Time 提取本地日期分量（时分秒归零）。
func LocalDateOf(v time.Time) LocalDate {
	dateOnly := time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, time.Local)
	return LocalDate{dateOnly}
}

// LocalDatePOf 等同 LocalDateOf，返回指针。
func LocalDatePOf(t time.Time) *LocalDate {
	dateOnly := LocalDateOf(t)
	return &dateOnly
}

// ----------------------- of Loc----------------------------

// LocalDateOfLoc 先将 v 转为 time.Local，再提取日期分量。
func LocalDateOfLoc(v time.Time) LocalDate {
	t := v.In(time.Local)
	dateOnly := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	return LocalDate{dateOnly}
}

// LocalDatePOfLoc 等同 LocalDateOfLoc，返回指针。
func LocalDatePOfLoc(t time.Time) *LocalDate {
	dateOnly := LocalDateOfLoc(t)
	return &dateOnly
}

// ----------------------- of Ymd ----------------------------

// LocalDateOfYmd 由年月日整数构造本地日期。
func LocalDateOfYmd(year, month, day int) LocalDate {
	dateOnly := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return LocalDate{dateOnly}
}

// LocalDatePOfYmd 等同 LocalDateOfYmd，返回指针。
func LocalDatePOfYmd(year, month, day int) *LocalDate {
	dateOnly := LocalDateOfYmd(year, month, day)
	return &dateOnly
}

// ----------------------- base ----------------------------

// String 返回 "2006-01-02" 格式字符串。
func (t LocalDate) String() string {
	return t.data.Format(`2006-01-02`)
}

// IsZero 判断是否为 0001-01-01。
func (t LocalDate) IsZero() bool {
	return t.String() == "0001-01-01"
}

// Add 在日历语义下加减日期；d 为 nil 时返回原值。
func (t LocalDate) Add(d *DurationOption) LocalDate {
	if d == nil {
		return t
	}
	return LocalDate{t.data.AddDate(d.year, d.month, d.day)}
}

// AddTime 将本地日期与墙钟时间组合为 LocalDateTime。
func (d LocalDate) AddTime(t LocalTime) LocalDateTime {
	return LocalDateTimeOfYmdHms(
		d.data.Year(), int(d.data.Month()), d.data.Day(),
		t.hour, t.minute, t.second,
	)
}

// ----------------------- comp ----------------------------

// Before 当 t 早于 d 时返回 true。
func (t LocalDate) Before(d LocalDate) bool {
	return t.ToGoTime().Before(d.ToGoTime())
}

// After 当 t 晚于 d 时返回 true。
func (t LocalDate) After(d LocalDate) bool {
	return t.ToGoTime().After(d.ToGoTime())
}

// Eq 当 t 与 d 表示同一天时返回 true。
func (t LocalDate) Eq(d LocalDate) bool {
	return t.ToGoTime() == d.ToGoTime()
}

// ----------------------- to ----------------------------

// ToGoTime 返回当天 00:00:00 的 time.Time（time.Local）。
func (t LocalDate) ToGoTime() time.Time {
	return time.Date(t.data.Year(), t.data.Month(), t.data.Day(), 0, 0, 0, 0, time.Local)
}

// ToDateTime 转为 LocalDateTime（时分秒为 00:00:00）。
func (t LocalDate) ToDateTime() LocalDateTime {
	return LocalDateTime{t.ToGoTime()}
}

// ToDateTimeP 等同 ToDateTime，返回指针。
func (t LocalDate) ToDateTimeP() *LocalDateTime {
	return &LocalDateTime{t.ToGoTime()}
}

// ----------------------- parse ----------------------------

// LocalDateParse 解析 "2006-01-02" 格式字符串。
func LocalDateParse(data string) (LocalDate, error) {
	localTime, err := time.ParseInLocation(`2006-01-02`, data, time.Local)
	return LocalDate{localTime}, err
}

// LocalDateParseMust 内部调用 LocalDateParse，失败时 panic。
func LocalDateParseMust(data string) LocalDate {
	localTime, err := LocalDateParse(data)
	if err != nil {
		panic(err)
	}
	return localTime
}

// LocalDateParseMustP 内部调用 LocalDateParse，失败时 panic，返回指针。
func LocalDateParseMustP(data string) *LocalDate {
	localTime, err := LocalDateParse(data)
	if err != nil {
		panic(err)
	}
	return &localTime
}

// ----------------------- json ----------------------------

// MarshalJSON 输出带引号的日期字符串。
func (t LocalDate) MarshalJSON() ([]byte, error) {
	return []byte(t.data.Format(`"2006-01-02"`)), nil
}

// UnmarshalJSON 解析 JSON 日期字符串；null 时不修改接收方。
func (t *LocalDate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	localTime, err := time.ParseInLocation(`"2006-01-02"`, string(data), time.Local)
	*t = LocalDate{data: localTime}
	return err
}

// ----------------------- db ----------------------------

// Value 写入数据库的 date 字符串（"2006-01-02"）。
func (t LocalDate) Value() (driver.Value, error) {
	return t.data.Format("2006-01-02"), nil
}

// Scan 从数据库读取；nil 时不修改接收方。支持 string、[]byte、time.Time、LocalDate。
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

// LocalDateList 映射 PostgreSQL date[]，实现 driver.Valuer 与 sql.Scanner。
type LocalDateList []LocalDate

// Value 返回 date[] 的 PostgreSQL 文本字面量。
func (p LocalDateList) Value() (driver.Value, error) {
	marshal, err := json.Marshal([]LocalDate(p))
	if err != nil {
		return nil, err
	}
	s := string(marshal)
	if s != "null" {
		s = "{" + s[1:len(s)-1] + "}"
	} else {
		s = "{}"
	}
	return s, nil
}

// Scan 从 PostgreSQL date[] 解析；nil 时不修改接收方。
func (p *LocalDateList) Scan(data any) error {
	if data == nil {
		return nil
	}
	var dates pgtype.FlatArray[pgtype.Date]
	if err := scanPgArray(pgtype.DateArrayOID, data, &dates); err != nil {
		return err
	}
	list := make([]LocalDate, len(dates))
	for i, element := range dates {
		if element.Valid {
			list[i] = LocalDateOf(element.Time)
		}
	}
	*p = list
	return nil
}
