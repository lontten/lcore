# lcore

[English](README.md) | [中文](README.zh.md)

Go utility library providing value types for **PostgreSQL + GORM** workflows: local date/time, UUID, decimal, and PostgreSQL array lists. Types implement `json.Marshaler`, `driver.Valuer`, and `sql.Scanner` for seamless API and database interoperability.

## Features

- **Local time** — `LocalDate`, `LocalTime`, `LocalDateTime`: wall-clock semantics (like Java `LocalDate`), always using `time.Local`
- **PostgreSQL arrays** — `StringList`, `BoolList`, `IntList`, `DecimalList`, `UUIDList`: map to `text[]`, `bool[]`, `uuid[]`, etc.
- **Mixed array** — `Array`: mixed-type PostgreSQL arrays for GORM custom fields
- **UUID** — `UUID`, `UUIDList`: JSON as 32-char hex without dashes; standard UUID strings in the database
- **Utilities** — `NullUint64`, `NilToZero`, `NewInt` / `NewString` / `NewBool`, `Fields`: pointer helpers, nil-to-zero, and more

## Installation

```bash
# v2 (current module, recommended)
go get -u github.com/lontten/lcore/v2

# v1 (legacy)
go get -u github.com/lontten/lcore
```

```go
import "github.com/lontten/lcore/v2/types"
```

## Quick Start

**Local date JSON:**

```go
d := types.LocalDateOfYmd(2026, 6, 23)
b, _ := json.Marshal(d) // "2026-06-23"
```

**UUID JSON (no dashes):**

```go
id := types.Str2UUIDMust("550e8400-e29b-41d4-a716-446655440000")
b, _ := json.Marshal(id) // "550e8400e29b41d4a716446655440000"
```

**PostgreSQL array field in a GORM model:**

```go
type Model struct {
    Tags types.StringList `gorm:"type:text[]"`
}
```

## Development

Requires Go 1.25+. This is a library module with no `main` package.

```bash
go mod verify
go test -race -count=1 ./...
```

## License

Apache 2.0 — see [LICENSE](LICENSE).
