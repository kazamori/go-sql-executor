package cmd_test

import (
	"context"
	"testing"

	"github.com/kazamori/go-sql-executor/cmd"
	"github.com/kazamori/go-sql-executor/concurrent"
	"github.com/kazamori/go-sql-executor/stats"
)

func TestFlatten(t *testing.T) {
	var data = []struct {
		name       string
		concurrent int
		values1    []float64
		values2    []float64
	}{
		{
			name:       "one",
			concurrent: 1,
			values1:    []float64{1.0, 2.0, 3.0},
			values2:    []float64{1.0, 2.0},
		},

		{
			name:       "four",
			concurrent: 4,
			values1:    []float64{1.0, 2.0},
			values2:    []float64{1.0},
		},

		{
			name:       "eight",
			concurrent: 8,
			values1:    []float64{1.0},
			values2:    []float64{1.0, 2.0},
		},
	}

	keyElapsed := "elapsed"
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

				m := make(map[string]stats.TimeValues)
				tv1 := *stats.NewTimeValues("msec")
				for _, v := range tt.values1 {
					tv1.Append(v)
				}
				m["tv1"] = tv1
				tv2 := *stats.NewTimeValues("msec")
				for _, v := range tt.values2 {
					tv2.Append(v)
				}
				m["tv2"] = tv2
				data[keyElapsed] = m
				return nil
			})

			actual := cmd.Flatten(results)
			t.Log(actual)
			if len(actual) != 2 {
				t.Errorf("got %v, want 2", len(actual))
			}

			tv1 := actual["tv1"]
			expectedTv1Length := len(tt.values1) * tt.concurrent
			if tv1.Len() != expectedTv1Length {
				t.Errorf("got %v, want %v", tv1.Len(), expectedTv1Length)
			}

			tv2 := actual["tv2"]
			expectedTv2Length := len(tt.values2) * tt.concurrent
			if tv2.Len() != expectedTv2Length {
				t.Errorf("got %v, want %v", tv2.Len(), expectedTv2Length)
			}
		})
	}
}
