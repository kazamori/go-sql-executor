package concurrent_test

import (
	"context"
	"testing"

	"github.com/kazamori/go-sql-executor/concurrent"
)

func TestCall(t *testing.T) {
	var data = []struct {
		name       string
		concurrent int
	}{
		{
			name:       "zero",
			concurrent: 0,
		},

		{
			name:       "one",
			concurrent: 1,
		},

		{
			name:       "eight",
			concurrent: 8,
		},
	}
	key := "value"
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			results := concurrent.Call(ctx, tt.concurrent, func(
				ctx context.Context, ch chan concurrent.Data,
			) error {
				data := concurrent.Data{}
				defer func() {
					ch <- data
				}()
				data[key] = 1
				return nil
			})

			if len(results) != tt.concurrent {
				t.Errorf("got %v, want %v", len(results), tt.concurrent)
			}

			for _, data := range results {
				actual := data[key].(int)
				if actual != 1 {
					t.Errorf("got %v, want 1", actual)
				}
			}
		})
	}
}
