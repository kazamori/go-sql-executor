package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/cybozu-go/well"
	"github.com/google/subcommands"
	"github.com/kazamori/go-sql-executor/cmd"
)

var (
	version = flag.Bool("version", false, "version")
)

var (
	revision  string
	buildTime string
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(cmd.GetInfoCommand(), "info")
	subcommands.Register(cmd.GetQueryCommand(), "query")
	flag.Parse()
	well.LogConfig{}.Apply()

	if *version {
		fmt.Printf("Build on %s from revision: %s\n", buildTime, revision)
		os.Exit(0)
	}

	exitStatus := subcommands.ExitSuccess
	well.Go(func(ctx context.Context) error {
		exitStatus = subcommands.Execute(ctx)
		return nil
	})
	well.Stop()
	well.Wait()
	os.Exit(int(exitStatus))
}
