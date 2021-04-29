package sqlite3

import (
	_ "github.com/xo/usql/drivers/sqlite3"
)

const (
	DriverName = "sqlite3"
)

type Sqlite3 struct {
	name string
}

func (s *Sqlite3) GetName() string {
	return s.name
}

func (s *Sqlite3) GetVersion() string {
	return "select sqlite_version()"
}

func New() *Sqlite3 {
	return &Sqlite3{
		name: DriverName,
	}
}
