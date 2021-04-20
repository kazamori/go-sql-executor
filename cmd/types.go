package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/kazamori/go-sql-executor/concurrent"
	"github.com/kazamori/go-sql-executor/db"
)

type envFlag struct {
	value string
}

func (f *envFlag) String() string {
	return "env"
}

func (f *envFlag) Set(value string) error {
	f.value = value
	return nil
}

var errUnsupportedDriver = errors.New("unsupported driver")

type driverFlag struct {
	driver db.Driver
}

func (f *driverFlag) String() string {
	return "driver"
}

func (f *driverFlag) Set(value string) error {
	driver, ok := db.GetDriver(value)
	if !ok {
		return errUnsupportedDriver
	}
	f.driver = driver
	return nil
}

type commonOption struct {
	driver driverFlag
	host   envFlag
	path   envFlag
	port   envFlag
	user   envFlag
	passwd envFlag
	schema envFlag
}

func newCommonOption() commonOption {
	driver, _ := db.GetDriver(os.Getenv("DB_DRIVER"))
	return commonOption{
		driver: driverFlag{driver: driver},
		host:   envFlag{value: os.Getenv("DB_HOST")},
		path:   envFlag{value: os.Getenv("DB_PATH")},
		port:   envFlag{value: os.Getenv("DB_PORT")},
		user:   envFlag{value: os.Getenv("DB_USER")},
		passwd: envFlag{value: os.Getenv("DB_PASSWORD")},
		schema: envFlag{value: os.Getenv("DB_SCHEMA")},
	}
}

type fileFlag struct {
	path  string
	lines []string
}

func (f *fileFlag) String() string {
	return "file"
}

func (f *fileFlag) Set(value string) error {
	r, err := os.Open(value)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", value, err)
	}

	lines, err := concurrent.ReadLines(r)
	if err != nil {
		return fmt.Errorf("failed to read lines: %w", err)
	}

	f.path = value
	f.lines = lines
	return nil
}
