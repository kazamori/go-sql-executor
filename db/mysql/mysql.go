package mysql

import (
	_ "github.com/xo/usql/drivers/mysql"
)

const (
	DriverName = "mysql"
)

type Mysql struct {
	name string
}

func (m *Mysql) GetName() string {
	return m.name
}

func (m *Mysql) GetVersion() string {
	return "select version()"
}

func New() *Mysql {
	return &Mysql{
		name: DriverName,
	}
}
