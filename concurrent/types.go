package concurrent

import (
	"context"
)

type Data map[string]interface{}

type Line struct {
	Value string
	Error error
}

type Func func(ctx context.Context, ch chan Data) error
