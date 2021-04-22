# go-sql-executor

A CLI tool to execute SQL queries for load-testing and performance analysis.

## How to build

```bash
$ make
go build -o bin/sql-executor -ldflags "-X main.revision=e7762737 -X main.buildTime=2021-04-21T12:43:14Z" cmd/executor/main.go
```

## Development Status

**Alpha.**

go-sql-executor is still under heavy development. Some functionality are known to be broken, missing or incomplete. The interface may also change.

## How to run

```bash
$ ./bin/sql-executor --help
Usage: sql-executor <flags> <subcommand> <subcommand args>

Subcommands:
    commands         list all command names
    flags            describe all known top-level flags
    help             describe subcommands and their syntax

Subcommands for info:
    info             show database information.

Subcommands for query:
    query            query any SQL.


Use "sql-executor flags" for a list of top-level flags
```

### Database configuration

The database configuration can be set via environmental variables.

* `DB_DRIVER`
* `DB_HOST`
* `DB_PATH`
* `DB_PORT`
* `DB_USER`
* `DB_PASSWORD`
* `DB_SCHEMA`

Also, you can overwrite these configuration by passing CLI options. CLI options take precedence over environmental variables.

```bash
  -driver value
        driver name (default from $DB_DRIVER)
  -host value
        host name (default from $DB_HOST)
  -path value
        path to dbfile (default from $DB_PATH)
  -user value
        db user (default from $DB_USER)
  -password value
        db password (default from $DB_PASSWORD)
  -port value
        port number (default from $DB_PORT)
  -schema value
        schema/dbname (default from $DB_SCHEMA)
```

For example, pass all options to connect the database.

```bash
$ ./bin/sql-executor query \
    -driver postgres \
    -host locahost \
    -port 5432 \
    -user postgres \
    -password secret \
    -schema test
```

To confirm the `driver` name supported by `sql-executor`, use use `info` subcommand.

```bash
$ ./bin/sql-executor info -drivers
available drivers:
  - mysql
  - postgres
  - sqlite3
```

### query subcommand

The `query` subcommand queries any SQL.

```bash
$ ./bin/sql-executor query --help
query:
  query any SQL.

  (omit database configuration options)

  -concurrent int
        the number of concurrent (default 1)
  -enableOutput
        output SQL results
  -enableTransaction
        execute as a transaction
  -file value
        a file including SQL queries
  -q string
        any SQL to query
  -repeat int
        repeat query given SQL (default 3)
```

For simple use, pass SQL with `-q` option.

```bash
$ ./bin/sql-executor query -q "select count(*) from mytable"
```

For complicated use case, pass a file including SQLs with `-file` option.

**NOTE: sql-executor runs a line (expects single SQL query) at a time**

```bash
$ vi test.sql
select count(distinct(name)) from mytable
select name from othertable
$ ./bin/sql-executor query -file test.sql
```

To confirm performance when multiple users access, use `-concurrent` option. In most cases, it results higher latency than querying solely.

```bash
$ ./bin/sql-executor query -file test.sql -concurrent 20
```

### info subcommand

The `info` subcommand shows database information. It is useful for a connection test or debugging.

```bash
$ ./bin/sql-executor info --help
info:
  show database information.

  (omit database configuration options)

  -drivers
        show available drivers
```

It shows the database version.

```bash
$ ./bin/sql-executor info
Connected with driver postgres (PostgreSQL 11.10)
  version
-----------------------------------------------------------------------------------------
  PostgreSQL 11.10 on x86_64-pc-linux-musl, compiled by gcc (Alpine 9.3.0) 9.3.0, 64-bit
(1 row)
```

