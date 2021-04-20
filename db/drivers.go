package db

import (
	"fmt"

	"github.com/kazamori/go-sql-executor/db/mysql"
	"github.com/kazamori/go-sql-executor/db/postgres"
	"github.com/kazamori/go-sql-executor/db/sqlite3"
)

var driverMap = map[string]Driver{
	mysql.DriverName:    mysql.New(),
	postgres.DriverName: postgres.New(),
	sqlite3.DriverName:  sqlite3.New(),
}

func GetDriver(name string) (driver Driver, ok bool) {
	driver, ok = driverMap[name]
	return
}

func GetAvailableDrivers() (names []string) {
	names = make([]string, 0, len(driverMap))
	for name := range driverMap {
		names = append(names, name)
	}
	return
}

type DataSourceConfig struct {
	Driver Driver
	Path   string
	Host   string
	Port   string
	User   string
	Passwd string
	Schema string
}

func GetDataSourceName(c *DataSourceConfig) string {
	driverName := c.Driver.GetName()
	switch driverName {
	case sqlite3.DriverName:
		return fmt.Sprintf("%s://%s", driverName, c.Path)
	default:
		if c.Passwd == "" {
			return fmt.Sprintf(
				"%s://%s@%s:%s/%s",
				driverName, c.User, c.Host, c.Port, c.Schema,
			)
		}
		return fmt.Sprintf(
			"%s://%s:%s@%s:%s/%s",
			driverName, c.User, c.Passwd, c.Host, c.Port, c.Schema,
		)
	}
}

func NewDataSourceConfig(
	driver Driver,
	host, path, port, user, passwd, schema string,
) *DataSourceConfig {
	return &DataSourceConfig{
		Driver: driver,
		Path:   path,
		Host:   host,
		Port:   port,
		User:   user,
		Passwd: passwd,
		Schema: schema,
	}
}
