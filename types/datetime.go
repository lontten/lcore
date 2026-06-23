package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// LocalDateTime 本地日期时间（墙上时钟，含年月日时分秒，不含时区偏移）。
// 语义等同 java.time.LocalDateTime，统一使用 time.Local。
// 零值为 0001-01-01 00:00:00。
type LocalDateTime struct {
	data time.Time
}

// ----------------------- now ----------------------------

// NowDateTime 返回当前本地日期时间。
func NowDateTime() LocalDateTime {
	return LocalDateTime{time.Now()}
}

// NowDateTimeP 返回当前本地日期时间的指针。
func NowDateTimeP() *LocalDateTime {
	return &LocalDateTime{time.Now()}
}

// ----------------------- of ----------------------------

// LocalDateTimeOf 从 time.Time 提取本地日期时间分量（纳秒归零）。
func LocalDateTimeOf(v time.Time) LocalDateTime {
	dateOnly := time.Date(v.Year(), v.Month(), v.Day(), v.Hour(), v.Minute(), v.Second(), 0, time.Local)
	return LocalDateTime{dateOnly}
}

// LocalDateTimePOf 等同 LocalDateTimeOf，返回指针。
func LocalDateTimePOf(v time.Time) *LocalDateTime {
	dateOnly := LocalDateTimeOf(v)
	return &dateOnly
}

// ----------------------- of Loc----------------------------

// LocalDateTimeOfLoc 先将 v 转为 time.Local，再作为本地日期时间。
func LocalDateTimeOfLoc(v time.Time) LocalDateTime {
	t := v.In(time.Local)
	return LocalDateTime{t}
}

// LocalDateTimePOfLoc 等同 LocalDateTimeOfLoc，返回指针。
func LocalDateTimePOfLoc(v time.Time) *LocalDateTime {
	dateOnly := LocalDateTimeOfLoc(v)
	return &dateOnly
}

// ----------------------- of YmdHms----------------------------

// LocalDateTimeOfYmdHms 由年月日时分秒整数构造本地日期时间。
func LocalDateTimeOfYmdHms(year, month, day, hour, min, sec int) LocalDateTime {
	dateOnly := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local)
	return LocalDateTime{dateOnly}
}

// LocalDateTimePOfYmdHms 等同 LocalDateTimeOfYmdHms，返回指针。
func LocalDateTimePOfYmdHms(year, month, day, hour, min, sec int) *LocalDateTime {
	dateOnly := LocalDateTimeOfYmdHms(year, month, day, hour, min, sec)
	return &dateOnly
}

// ----------------------- base ----------------------------

// String 返回 "2006-01-02 15:04:05" 格式字符串。
func (t LocalDateTime) String() string {
	return t.data.Format(`2006-01-02 15:04:05`)
}

// IsZero 判断是否为 0001-01-01 00:00:00。
func (t LocalDateTime) IsZero() bool {
	return t.String() == "0001-01-01 00:00:00"
}

// Add 在本地日期时间语义下加减时长；d 为 nil 时返回原值。
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

// Before 当 t 早于 d 时返回 true。
func (t LocalDateTime) Before(d LocalDateTime) bool {
	return t.ToGoTime().Before(d.ToGoTime())
}

// After 当 t 晚于 d 时返回 true。
func (t LocalDateTime) After(d LocalDateTime) bool {
	return t.ToGoTime().After(d.ToGoTime())
}

// Eq 当 t 与 d 相等时返回 true。
func (t LocalDateTime) Eq(d LocalDateTime) bool {
	return t.ToGoTime() == d.ToGoTime()
}

// ----------------------- to ----------------------------

// ToGoTime 返回底层 time.Time。
func (t LocalDateTime) ToGoTime() time.Time {
	return t.data
}

// ToDate 提取日期部分为 LocalDate。
func (t LocalDateTime) ToDate() LocalDate {
	return LocalDate{t.ToGoTime()}
}

// ToDateP 等同 ToDate，返回指针。
func (t LocalDateTime) ToDateP() *LocalDate {
	return &LocalDate{t.ToGoTime()}
}

// ----------------------- parse ----------------------------

// LocalDateTimeParse 解析 "2006-01-02 15:04:05" 格式字符串。
func LocalDateTimeParse(data string) (LocalDateTime, error) {
	localTime, err := time.ParseInLocation(`2006-01-02 15:04:05`, data, time.Local)
	return LocalDateTime{localTime}, err
}

// LocalDateTimeParseMust 内部调用 LocalDateTimeParse，失败时 panic。
func LocalDateTimeParseMust(data string) LocalDateTime {
	localTime, err := LocalDateTimeParse(data)
	if err != nil {
		panic(err)
	}
	return localTime
}

// LocalDateTimeParseMustP 内部调用 LocalDateTimeParse，失败时 panic，返回指针。
func LocalDateTimeParseMustP(data string) *LocalDateTime {
	localTime, err := LocalDateTimeParse(data)
	if err != nil {
		panic(err)
	}
	return &localTime
}

// ----------------------- json ----------------------------

// MarshalJSON 输出带引号的日期时间字符串。
func (t LocalDateTime) MarshalJSON() ([]byte, error) {
	return []byte(t.data.Format(`"2006-01-02 15:04:05"`)), nil
}

// UnmarshalJSON 解析 JSON 日期时间字符串；null 时不修改接收方。
func (t *LocalDateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	localTime, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	*t = LocalDateTime{data: localTime}
	return err
}

// ----------------------- db ----------------------------

// Value 写入数据库的 timestamp（返回 time.Time，供驱动按 timestamp 列处理）。
func (t LocalDateTime) Value() (driver.Value, error) {
	return t.data, nil
}

// Scan 从数据库读取；nil 时不修改接收方。支持 string、[]byte、time.Time、LocalDateTime。
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
		return nil
	case LocalDateTime:
		*t = v
		return nil
	default:
		return fmt.Errorf("can not convert %v to LocalDateTime", v)
	}
	if len(s) < 19 {
		return fmt.Errorf("can not convert %v to LocalDateTime", v)
	}
	localTime, err := time.ParseInLocation(`2006-01-02 15:04:05`, s, time.Local)
	if err != nil {
		return err
	}
	*t = LocalDateTime{data: localTime}
	return nil
}

// ----------------------- list ----------------------------

// LocalDateTimeList 映射 PostgreSQL timestamp[]，实现 driver.Valuer 与 sql.Scanner。
type LocalDateTimeList []LocalDateTime

// Value 返回 timestamp[] 的 PostgreSQL 文本字面量。
func (p LocalDateTimeList) Value() (driver.Value, error) {
	marshal, err := json.Marshal([]LocalDateTime(p))
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

// Scan 从 PostgreSQL timestamp[] 解析；nil 时不修改接收方。
func (p *LocalDateTimeList) Scan(data any) error {
	if data == nil {
		return nil
	}
	var stamps pgtype.FlatArray[pgtype.Timestamp]
	if err := scanPgArray(pgtype.TimestampArrayOID, data, &stamps); err != nil {
		return err
	}
	list := make([]LocalDateTime, len(stamps))
	for i, element := range stamps {
		if element.Valid {
			list[i] = LocalDateTimeOf(element.Time)
		}
	}
	*p = list
	return nil
}
