package query

import (
	"context"
	_sql "database/sql"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/cybozu-go/log"
	"github.com/kazamori/go-sql-executor/db"
	"github.com/kazamori/go-sql-executor/stats"
	"github.com/xo/usql/handler"
	"github.com/xo/usql/rline"
)

type Handler struct {
	config       *db.DataSourceConfig
	enableOutput bool

	raw         *handler.Handler
	dsn         string
	elapsedTime map[string]stats.TimeValues
}

func (h *Handler) configure() (err error) {
	user_, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	var line rline.IO
	if h.enableOutput {
		line, err = rline.New(true, "", "")
	} else {
		line, err = NewNopRline()
	}
	if err != nil {
		return fmt.Errorf("failed to get readline: %w", err)
	}

	h.raw = handler.New(line, user_, wd, true)
	h.raw.SetSingleLineMode(true)
	return nil
}

func (h *Handler) Connect() error {
	if err := h.configure(); err != nil {
		return fmt.Errorf("failed to configure: %w", err)
	}
	return h.raw.Open(context.TODO(), h.dsn)
}

func (h *Handler) setElapsedTime(sql string, msec float64) {
	tv, ok := h.elapsedTime[sql]
	if !ok {
		tv = *stats.NewTimeValues("msec")
	}
	tv.Append(msec)
	h.elapsedTime[sql] = tv
}

func (h *Handler) Query(
	ctx context.Context, rawSql string, useTransaction bool,
) (err error) {
	var sql = rawSql
	if HasBindVariable(rawSql) {
		p := NewSimpleParser(rawSql)
		p.Parse()
		sql, err = p.Replace()
		if err != nil {
			return err
		}
	}

	if useTransaction {
		var txOpts = &_sql.TxOptions{
			Isolation: _sql.LevelReadUncommitted,
			ReadOnly:  false,
		}
		if err = h.raw.Begin(txOpts); err != nil {
			return err
		}
	}

	h.raw.Reset([]rune(sql))
	start := time.Now()
	err = h.raw.Run()

	if useTransaction {
		if e := h.raw.Commit(); e != nil {
			return e
		}
	}

	elapsed := time.Since(start)
	log.Debug("elapsed time", map[string]interface{}{
		"raw_sql":  rawSql,
		"exec_sql": sql,
		"took":     elapsed,
	})
	mseconds := float64(elapsed.Microseconds()) / 1000
	h.setElapsedTime(sql, mseconds)
	return err
}

func (h *Handler) ShowSystemInfo() error {
	return h.Query(
		context.Background(), h.config.Driver.GetVersion(), false,
	)
}

func (h *Handler) GetElapsedTime() map[string]stats.TimeValues {
	return h.elapsedTime
}

func NewHandler(
	c *db.DataSourceConfig, enableOutput bool,
) *Handler {
	return &Handler{
		config:       c,
		enableOutput: enableOutput,
		dsn:          db.GetDataSourceName(c),
		elapsedTime:  make(map[string]stats.TimeValues),
	}
}
