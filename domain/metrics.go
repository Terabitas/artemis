package domain

import "time"

type (
	// Metric type
	Metric interface {
		GetValue() float64
		GetTimestamp() time.Time
	}

	// HealthMetric metric
	HealthMetric struct {
		Value float64
		Time  time.Time
	}

	// MetricSeries type
	MetricSeries map[time.Time]Metric
)

// NewHealthMetric constructor
func NewMetricSeries(m ...Metric) MetricSeries {
	ms := MetricSeries{}
	for _, mToAdd := range m {
		ms[mToAdd.GetTimestamp()] = mToAdd
	}

	return ms
}

// NewHealthMetric constructor
func NewHealthMetric(val float64, t time.Time) Metric {
	return HealthMetric{
		Value: val,
		Time:  t,
	}
}

// GetValue ...
func (hm HealthMetric) GetValue() float64 {
	return hm.Value
}

// GetTimestamp ...
func (hm HealthMetric) GetTimestamp() time.Time {
	return hm.Time
}
