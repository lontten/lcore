// Package types 提供 PostgreSQL/GORM 场景下的本地时间、数组列表、UUID、decimal 等值类型与工具函数。
//
// 时间类型（LocalDate、LocalTime、LocalDateTime）表示墙上时钟语义，统一使用 time.Local，
// 并实现 json.Marshaler、driver.Valuer、sql.Scanner 以便与 API 和数据库互通。
//
// *List 类型（如 StringList、BoolList）映射 PostgreSQL 数组列，供 GORM 自定义字段使用。
package types
