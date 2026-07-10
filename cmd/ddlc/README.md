# ddlc CLI
<img src="docs/img/badges.svg">

`ddlc` is a command-line tool that scans a codebase for model definitions and exports their database schema definitions (DDL) to a file or standard output without requiring a running database server.

## Installation

```bash
go install github.com/tinywasm/ddlc/cmd/ddlc@latest
```

## Usage

```bash
ddlc [flags]
```

### Flags

- **`-root`** (string): The root directory to scan recursively for `model.go` or `models.go` definition files. Defaults to the current directory (`.`).
- **`-out`** (string): The path to write the generated DDL schema. Use `-` to print directly to standard output. Defaults to `-`.
- **`-dialect`** (string): The target SQL database dialect: `sqlite` or `postgres`. Defaults to `sqlite`.

## Examples

### Output SQLite DDL to standard output

```bash
ddlc -dialect sqlite
```

### Save PostgreSQL schema DDL to a file

```bash
ddlc -root ./pkg/models -dialect postgres -out schema.sql
```
