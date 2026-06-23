# lcore

[English](README.md) | [中文](README.zh.md)

Go 工具库，为 **PostgreSQL + GORM** 场景提供本地时间、UUID、decimal、PG 数组列表等值类型；实现 `json.Marshaler`、`driver.Valuer`、`sql.Scanner`，便于 API 与数据库互通。

## 功能特性

- **本地时间** — `LocalDate`、`LocalTime`、`LocalDateTime`：墙上时钟语义（等同 Java `LocalDate` 等），统一使用 `time.Local`
- **PG 数组** — `StringList`、`BoolList`、`IntList`、`DecimalList`、`UUIDList`：映射 `text[]`、`bool[]`、`uuid[]` 等
- **混合数组** — `Array`：GORM 自定义字段用混合类型 PG 数组
- **UUID** — `UUID`、`UUIDList`：JSON 无连字符 32 位 hex；数据库为标准 UUID 字符串
- **工具** — `NullUint64`、`NilToZero`、`NewInt` / `NewString` / `NewBool`、`Fields`：指针构造、nil 转零值等

## 安装

```bash
# v2（当前模块，推荐）
go get -u github.com/lontten/lcore/v2

# v1（历史版本）
go get -u github.com/lontten/lcore
```

```go
import "github.com/lontten/lcore/v2/types"
```

## 快速上手

**本地日期 JSON：**

```go
d := types.LocalDateOfYmd(2026, 6, 23)
b, _ := json.Marshal(d) // "2026-06-23"
```

**UUID JSON（无连字符）：**

```go
id := types.Str2UUIDMust("550e8400-e29b-41d4-a716-446655440000")
b, _ := json.Marshal(id) // "550e8400e29b41d4a716446655440000"
```

**GORM 模型中的 PG 数组字段：**

```go
type Model struct {
    Tags types.StringList `gorm:"type:text[]"`
}
```

## 开发与测试

需要 Go 1.25+。本仓库为库模块，无 `main` 包。

```bash
go mod verify
go test -race -count=1 ./...
```

## 许可证

Apache 2.0 — 见 [LICENSE](LICENSE)。
