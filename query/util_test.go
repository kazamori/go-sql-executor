package query_test

import (
	"reflect"
	"testing"

	"github.com/kazamori/go-sql-executor/query"
)

func TestZip(t *testing.T) {
	var data = []struct {
		name     string
		values1  []string
		values2  []string
		expected [][]string
	}{
		{
			name:    "simple",
			values1: []string{"a", "b", "c"},
			values2: []string{"x", "y", "z"},
			expected: [][]string{
				{"a", "x"},
				{"b", "y"},
				{"c", "z"},
			},
		},

		{
			name:    "unmatch",
			values1: []string{"a", "b"},
			values2: []string{"x", "y", "z"},
			expected: [][]string{
				{"a", "x"},
				{"b", "y"},
			},
		},
	}

	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			actual := query.Zip(tt.values1, tt.values2)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("got %v, want %v", actual, tt.expected)
			}
		})
	}
}
