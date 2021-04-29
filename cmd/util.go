package cmd

import (
	"fmt"

	"github.com/kazamori/go-sql-executor/concurrent"
	"github.com/kazamori/go-sql-executor/stats"
)

func validateCommonOption(
	c commonOption,
) bool {
	if c.driver.driver == nil {
		fmt.Println("required -driver argument")
		return false
	}
	if c.host.value == "" {
		fmt.Println("required -host argument")
		return false
	}
	if c.port.value == "" {
		fmt.Println("required -port argument")
		return false
	}
	if c.user.value == "" {
		fmt.Println("required -user argument")
		return false
	}
	if c.schema.value == "" {
		fmt.Println("required -schema argument")
		return false
	}
	return true
}

func Flatten(results []concurrent.Data) map[string]stats.TimeValues {
	flattened := make(map[string]stats.TimeValues)
	for _, data := range results {
		if elapsed, exist := data[keyElapsed]; exist {
			if m, ok := elapsed.(map[string]stats.TimeValues); ok {
				for key, _tv := range m {
					tv, has := flattened[key]
					if has {
						tv.AppendTimeValue(_tv)
					} else {
						tv = _tv
					}
					flattened[key] = tv
				}
			}
		}
	}
	return flattened
}
