package stats

type TimeValues struct {
	values []float64
	unit   string
}

func (e *TimeValues) Append(value float64) {
	e.values = append(e.values, value)
}

func (e *TimeValues) AppendTimeValue(tv TimeValues) {
	for _, value := range tv.values {
		e.Append(value)
	}
}

func NewTimeValues(unit string) *TimeValues {
	return &TimeValues{
		values: make([]float64, 0),
		unit:   unit,
	}
}
