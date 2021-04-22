package stats

import (
	"fmt"

	"github.com/montanaflynn/stats"
)

type Stats struct {
	min     float64
	max     float64
	mean    float64
	median  float64
	stddevp float64
	stddevs float64

	unit string
}

func (s *Stats) Show() {
	fmt.Printf("basic (%s):\n", s.unit)
	fmt.Printf(" - min    : %.3f\n", s.min)
	fmt.Printf(" - max    : %.3f\n", s.max)
	fmt.Printf(" - mean   : %.3f\n", s.mean)
	fmt.Printf(" - median : %.3f\n", s.median)
	fmt.Printf(" - stddevp: %.3f\n", s.stddevp)
	fmt.Printf(" - stddevs: %.3f\n", s.stddevs)
}

func GetBasicStatistics(values []float64, unit string) (*Stats, error) {
	min, err := stats.Min(values)
	if err != nil {
		return nil, fmt.Errorf("failed to get Min: %w", err)
	}
	max, err := stats.Max(values)
	if err != nil {
		return nil, fmt.Errorf("failed to get Max: %w", err)
	}
	mean, err := stats.Mean(values)
	if err != nil {
		return nil, fmt.Errorf("failed to get Mean: %w", err)
	}
	median, _ := stats.Median(values)
	if err != nil {
		return nil, fmt.Errorf("failed to get Median: %w", err)
	}
	stddevp, _ := stats.StdDevP(values)
	if err != nil {
		return nil, fmt.Errorf("failed to get StdDevP: %w", err)
	}
	stddevs, _ := stats.StdDevS(values)
	if err != nil {
		return nil, fmt.Errorf("failed to get StdDevS: %w", err)
	}

	stats := &Stats{
		min:     min,
		max:     max,
		mean:    mean,
		median:  median,
		stddevp: stddevp,
		stddevs: stddevs,
		unit:    unit,
	}
	return stats, nil
}

type Percentile struct {
	percent float64
	value   float64
}

func (p *Percentile) Show() {
	fmt.Printf(" - p%d: %.3f\n", int(p.percent), p.value)
}

type Percentiles struct {
	values []Percentile

	unit string
}

func (p *Percentiles) Append(value Percentile) {
	p.values = append(p.values, value)
}

func (p *Percentiles) Show() {
	fmt.Printf("percentiles (%s):\n", p.unit)
	for _, percentile := range p.values {
		percentile.Show()
	}
}

func GetPercentiles(
	values []float64, percents []float64, unit string,
) (*Percentiles, error) {
	percentiles := &Percentiles{
		values: make([]Percentile, 0),
		unit:   unit,
	}
	for _, percent := range percents {
		value, err := stats.Percentile(values, percent)
		if err != nil {
			return nil, err
		}

		percentiles.Append(Percentile{
			percent: percent,
			value:   value,
		})
	}
	return percentiles, nil
}

var (
	defaultPercents = []float64{
		99, 95, 90, 80, 70,
	}
)

func ShowStatistics(data map[string]TimeValues) error {
	for key, tv := range data {
		fmt.Printf("target:\n%s\n\n", key)

		basic, err := GetBasicStatistics(tv.values, tv.unit)
		if err != nil {
			return fmt.Errorf("failed to get basics: %w", err)
		}
		basic.Show()

		fmt.Println()

		percentiles, err := GetPercentiles(tv.values, defaultPercents, tv.unit)
		if err != nil {
			return fmt.Errorf("failed to get percentiles: %w", err)
		}
		percentiles.Show()

		fmt.Println()
	}
	return nil
}
