package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/cybozu-go/log"
	"github.com/google/subcommands"
	"github.com/kazamori/go-sql-executor/concurrent"
	"github.com/kazamori/go-sql-executor/db"
	"github.com/kazamori/go-sql-executor/query"
	"github.com/kazamori/go-sql-executor/stats"
)

type queryCmd struct {
	commonOption
	q                 string
	file              fileFlag
	concurrent        int
	repeat            int
	enableOutput      bool
	enableTransaction bool
}

func (*queryCmd) Name() string {
	return "query"
}

func (*queryCmd) Synopsis() string {
	return "query any SQL."
}

func (*queryCmd) Usage() string {
	return `query:
  query any SQL.
`
}

func (c *queryCmd) SetFlags(f *flag.FlagSet) {
	// common
	f.Var(&c.driver, "driver", "driver name (default from $DB_DRIVER)")
	f.Var(&c.host, "host", "host name (default from $DB_HOST)")
	f.Var(&c.path, "path", "path to dbfile (default from $DB_PATH)")
	f.Var(&c.port, "port", "port number (default from $DB_PORT)")
	f.Var(&c.user, "user", "db user (default from $DB_USER)")
	f.Var(&c.passwd, "password", "db password (default from $DB_PASSWORD)")
	f.Var(&c.schema, "schema", "schema/dbname (default from $DB_SCHEMA)")
	// query specific
	f.StringVar(&c.q, "q", "", "any SQL to query")
	f.Var(&c.file, "file", "a file including SQL queries")
	f.IntVar(&c.concurrent, "concurrent", 1, "the number of concurrent")
	f.IntVar(&c.repeat, "repeat", 3, "repeat query given SQL")
	f.BoolVar(&c.enableOutput, "enableOutput", false, "output SQL results")
	f.BoolVar(&c.enableTransaction,
		"enableTransaction", false, "execute as a transaction")
}

func (c *queryCmd) executeQuery(
	ctx context.Context, q string, queries []string, h *query.Handler,
) error {
	i := 0
	for i < c.repeat {
		if q != "" {
			if err := h.Query(ctx, q, c.enableTransaction); err != nil {
				return err
			}
		}

		if len(queries) > 0 {
			for _, q := range queries {
				if err := h.Query(ctx, q, c.enableTransaction); err != nil {
					return err
				}
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		i++
	}

	return nil
}

func (c *queryCmd) executeSingle(
	ctx context.Context, config *db.DataSourceConfig,
) subcommands.ExitStatus {
	h := query.NewHandler(config, c.enableOutput)
	if err := h.Connect(); err != nil {
		log.Error("failed to connect", map[string]interface{}{
			"schema": c.schema,
			"err":    err,
		})
		return subcommands.ExitFailure
	}

	if err := c.executeQuery(ctx, c.q, c.file.lines, h); err != nil {
		log.Error("failed to query", map[string]interface{}{
			"q":    c.q,
			"file": c.file.path,
			"err":  err,
		})
		return subcommands.ExitFailure
	}

	if err := stats.ShowStatistics(h.GetElapsedTime()); err != nil {
		log.Error("failed to get statistics", map[string]interface{}{
			"err": err,
		})
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

const (
	keyElapsed = "elapsed"
)

func (c *queryCmd) executeConcurrently(
	ctx context.Context, config *db.DataSourceConfig,
) subcommands.ExitStatus {
	results := concurrent.Call(ctx, c.concurrent, func(
		ctx context.Context, ch chan concurrent.Data,
	) (err error) {
		data := concurrent.Data{}
		defer func() {
			ch <- data
		}()
		h := query.NewHandler(config, c.enableOutput)
		if err = h.Connect(); err != nil {
			return err
		}
		if err = c.executeQuery(ctx, c.q, c.file.lines, h); err != nil {
			return err
		}
		data[keyElapsed] = h.GetElapsedTime()
		return nil
	})
	elapsed := Flatten(results)
	stats.ShowStatistics(elapsed)
	return subcommands.ExitSuccess
}

func (c *queryCmd) Execute(
	ctx context.Context, f *flag.FlagSet, _ ...interface{},
) subcommands.ExitStatus {
	if !validateCommonOption(c.commonOption) {
		f.Usage()
		return subcommands.ExitUsageError
	}
	if c.q == "" && c.file.path == "" {
		fmt.Println("required -q or -file argument")
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
	if c.concurrent == 1 {
		return c.executeSingle(ctx, config)
	}
	return c.executeConcurrently(ctx, config)
}

func GetQueryCommand() subcommands.Command {
	return &queryCmd{
		commonOption: newCommonOption(),
		file:         fileFlag{},
	}
}
