package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
)

// LocalTime 本地墙钟时间（仅时分秒，不含日期与时区偏移）。
// 语义等同 java.time.LocalTime，始终按 time.Local 理解，不做跨时区 instant 换算。
// 零值为 00:00:00；结构体零值 {} 与 LocalTimeOfHms(0,0,0) 在 == 下相等。
type LocalTime struct {
	hour, minute, second int // 分量存储，保证零值 canonical
}

// localTimeAnchorYear 为 toTime 使用的虚拟日期，仅用于时分秒运算。
const localTimeAnchorYear = 1

// toTime 将墙钟时间映射为锚定日期的 time.Time，供 Add 等运算复用标准库。
func (t LocalTime) toTime() time.Time {
	return time.Date(localTimeAnchorYear, time.January, 1,
		t.hour, t.minute, t.second, 0, time.Local)
}

// localTimeFrom 从 time.Time 提取本地墙钟时分秒（先转 time.Local）。
func localTimeFrom(tm time.Time) LocalTime {
	tm = tm.In(time.Local)
	return LocalTime{tm.Hour(), tm.Minute(), tm.Second()}
}

// localTimeFromHms 校验并构造 LocalTime，拒绝越界分量（如 25:00:00）。
func localTimeFromHms(h, m, s int) (LocalTime, error) {
	if h < 0 || h > 23 || m < 0 || m > 59 || s < 0 || s > 59 {
		return LocalTime{}, fmt.Errorf("invalid local time %02d:%02d:%02d", h, m, s)
	}
	return LocalTime{hour: h, minute: m, second: s}, nil
}

// mustLocalTimeFromHms 等同 localTimeFromHms，失败时 panic（供 OfHms 等 Must 风格 API 使用）。
func mustLocalTimeFromHms(h, m, s int) LocalTime {
	t, err := localTimeFromHms(h, m, s)
	if err != nil {
		panic(err)
	}
	return t
}

// ----------------------- now ----------------------------

// NowTime 返回当前本地墙钟时间。
func NowTime() LocalTime {
	return localTimeFrom(time.Now())
}

func NowTimeP() *LocalTime {
	t := NowTime()
	return &t
}

// ----------------------- of ----------------------------

// LocalTimeOf 直接取 v 的时分秒字段，不做时区转换。
// 适用于 v 已是本地墙钟语义（如从 DB 读出、time.Now()）。
func LocalTimeOf(v time.Time) LocalTime {
	return LocalTime{v.Hour(), v.Minute(), v.Second()}
}

func LocalTimePOf(t time.Time) *LocalTime {
	timeOnly := LocalTimeOf(t)
	return &timeOnly
}

// ----------------------- of Loc----------------------------

// LocalTimeOfLoc 先将 v 转为 time.Local，再取时分秒。
// 适用于 v 可能带 UTC 等非本地 location 的场景。
func LocalTimeOfLoc(v time.Time) LocalTime {
	return localTimeFrom(v)
}

func LocalTimePOfLoc(t time.Time) *LocalTime {
	timeOnly := LocalTimeOfLoc(t)
	return &timeOnly
}

// ----------------------- of Hms----------------------------

// LocalTimeOfHms 由整数构造墙钟时间，非法分量 panic。
func LocalTimeOfHms(hour, min, sec int) LocalTime {
	return mustLocalTimeFromHms(hour, min, sec)
}

func LocalTimePOfHms(hour, min, sec int) *LocalTime {
	timeOnly := LocalTimeOfHms(hour, min, sec)
	return &timeOnly
}

// ----------------------- base ----------------------------

func (t LocalTime) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.hour, t.minute, t.second)
}

func (t LocalTime) IsZero() bool {
	return t.hour == 0 && t.minute == 0 && t.second == 0
}

// Add 在墙钟语义下做时长加减；跨日进位折叠回 00:00:00~23:59:59（如 23:30 + 1h = 00:30）。
func (t LocalTime) Add(d *DurationOption) LocalTime {
	if d == nil {
		return t
	}
	return localTimeFrom(t.toTime().Add(
		time.Duration(d.hour)*time.Hour +
			time.Duration(d.min)*time.Minute +
			time.Duration(d.sec)*time.Second +
			time.Duration(d.nsec)*time.Nanosecond,
	))
}

// AddData 将本地日期与墙钟时间组合为 LocalDateTime。
func (t LocalTime) AddData(d LocalDate) LocalDateTime {
	return LocalDateTimeOfYmdHms(
		d.data.Year(), int(d.data.Month()), d.data.Day(),
		t.hour, t.minute, t.second,
	)
}

// ----------------------- comp ----------------------------

// Before
// t<d 返回true
// t>=d 返回false
func (t LocalTime) Before(d LocalTime) bool {
	if t.hour != d.hour {
		return t.hour < d.hour
	}
	if t.minute != d.minute {
		return t.minute < d.minute
	}
	return t.second < d.second
}

// After
// t>d 返回true
// t<=d 返回false
func (t LocalTime) After(d LocalTime) bool {
	if t.hour != d.hour {
		return t.hour > d.hour
	}
	if t.minute != d.minute {
		return t.minute > d.minute
	}
	return t.second > d.second
}

// Eq
// t==d 返回true
// t!=d 返回false
func (t LocalTime) Eq(d LocalTime) bool {
	return t.hour == d.hour && t.minute == d.minute && t.second == d.second
}

// ----------------------- to ----------------------------

// ToGoTime 返回锚定日期 0001-01-01 上的 time.Time，location 为 time.Local。
func (t LocalTime) ToGoTime() time.Time {
	return t.toTime()
}

func (t LocalTime) ToLocalDateTime() LocalDateTime {
	return LocalDateTime{t.toTime()}
}

func (t LocalTime) ToLocalDateTimeP() *LocalDateTime {
	dt := t.ToLocalDateTime()
	return &dt
}

// ----------------------- parse ----------------------------

// LocalTimeParse 解析 "15:04:05" 格式字符串，要求与规范输出完全一致（含前导零）。
func LocalTimeParse(data string) (LocalTime, error) {
	parsed, err := time.ParseInLocation(`15:04:05`, data, time.Local)
	if err != nil {
		return LocalTime{}, err
	}
	t := localTimeFrom(parsed)
	if t.String() != data {
		return LocalTime{}, fmt.Errorf("invalid local time %q", data)
	}
	return t, nil
}

func LocalTimeParseP(data string) (*LocalTime, error) {
	t, err := LocalTimeParse(data)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func LocalTimeParseMust(data string) LocalTime {
	localTime, err := LocalTimeParse(data)
	if err != nil {
		panic(err)
	}
	return localTime
}

func LocalTimeParseMustP(data string) *LocalTime {
	localTime, err := LocalTimeParseP(data)
	if err != nil {
		panic(err)
	}
	return localTime
}

// ----------------------- json ----------------------------

func (t LocalTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

// UnmarshalJSON 解析 JSON 字符串；null 时不修改接收方（保持原值）。
func (t *LocalTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	parsed, err := LocalTimeParse(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}
	*t = parsed
	return nil
}

// ----------------------- db ----------------------------

// Value 写入数据库的 time 字符串（"15:04:05"）。
func (t LocalTime) Value() (driver.Value, error) {
	return t.String(), nil
}

// Scan 从数据库读取；nil 时不修改接收方。字符串支持纯时间或带日期前缀（取末 8 位）。
func (t *LocalTime) Scan(v any) error {
	if v == nil {
		return nil
	}
	switch v := v.(type) {
	case string:
		parsed, err := localTimeScanString(v)
		if err != nil {
			return err
		}
		*t = parsed
		return nil
	case []byte:
		parsed, err := localTimeScanString(string(v))
		if err != nil {
			return err
		}
		*t = parsed
		return nil
	case time.Time:
		*t = LocalTimeOf(v)
		return nil
	case LocalTime:
		*t = v
		return nil
	default:
		return fmt.Errorf("can not convert %v to LocalTime", v)
	}
}

func localTimeScanString(s string) (LocalTime, error) {
	if len(s) < 8 {
		return LocalTime{}, fmt.Errorf("can not convert %q to LocalTime", s)
	}
	return LocalTimeParse(s[len(s)-8:])
}

// ----------------------- list ----------------------------

// LocalTimeList PostgreSQL time 数组的 GORM 自定义类型。
type LocalTimeList []LocalTime

// Value 实现 driver.Valuer
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

// Scan 实现 sql.Scanner
func (p *LocalTimeList) Scan(data any) error {
	array := pgtype.VarcharArray{}
	if err := array.Scan(data); err != nil {
		return err
	}
	list := make([]LocalTime, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = LocalTimeParseMust(element.String)
	}
	*p = list
	return nil
}
