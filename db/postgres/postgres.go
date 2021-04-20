package postgres

import (
	_ "github.com/xo/usql/drivers/postgres"
)

const (
	DriverName = "postgres"
)

type Postgres struct {
	name string
}

func (p *Postgres) GetName() string {
	return p.name
}

func (p *Postgres) GetVersion() string {
	return "select version()"
}

func New() *Postgres {
	return &Postgres{
		name: DriverName,
	}
}
