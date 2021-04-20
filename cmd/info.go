package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/cybozu-go/log"
	"github.com/google/subcommands"
	"github.com/kazamori/go-sql-executor/db"
	"github.com/kazamori/go-sql-executor/query"

	"github.com/xo/usql/drivers"
)

type infoCmd struct {
	drivers bool
	commonOption
}

func (*infoCmd) Name() string {
	return "info"
}

func (*infoCmd) Synopsis() string {
	return "show database information."
}

func (*infoCmd) Usage() string {
	return `info:
  show database information.
`
}

func (c *infoCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.drivers, "drivers", false, "show available drivers")
	f.Var(&c.driver, "driver", "driver name (default from $DB_DRIVER)")
	f.Var(&c.host, "host", "host name (default from $DB_HOST)")
	f.Var(&c.path, "path", "path to dbfile (default from $DB_PATH)")
	f.Var(&c.port, "port", "port number (default from $DB_PORT)")
	f.Var(&c.user, "user", "db user (default from $DB_USER)")
	f.Var(&c.passwd, "password", "db password (default from $DB_PASSWORD)")
	f.Var(&c.schema, "schema", "schema/dbname (default from $DB_SCHEMA)")
}

func (c *infoCmd) Execute(
	ctx context.Context, f *flag.FlagSet, _ ...interface{},
) subcommands.ExitStatus {
	if c.drivers {
		available := drivers.Available()
		fmt.Println("available drivers:")
		for _, name := range db.GetAvailableDrivers() {
			if _, ok := available[name]; ok {
				fmt.Printf("  - %s\n", name)
			}
		}
		return subcommands.ExitSuccess
	}

	if !validateCommonOption(c.commonOption) {
		f.Usage()
		return subcommands.ExitUsageError
	}

	config := db.NewDataSourceConfig(
		c.driver.driver,
		c.host.value,
		c.path.value,
		c.port.value,
		c.user.value,
		c.passwd.value,
		c.schema.value)
	h := query.NewHandler(config, true)
	if err := h.Connect(); err != nil {
		log.Error("failed to connect", map[string]interface{}{
			"schema": c.schema,
			"err":    err,
		})
		return subcommands.ExitFailure
	}

	if err := h.ShowSystemInfo(); err != nil {
		log.Error("failed to get system information", map[string]interface{}{
			"err": err,
		})
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func GetInfoCommand() subcommands.Command {
	return &infoCmd{
		commonOption: newCommonOption(),
	}
}
