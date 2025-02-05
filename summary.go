package zpm

import (
	"sync/atomic"

	dto "github.com/prometheus/client_model/go"
)

// Summary client API
type summary struct {
	name      string
	help      *string
	unit      *string
	labels    []*dto.LabelPair
	quantiles []float64
	storage   *storage
}

func (s *summary) Help(help string) *summary {
	s.help = &help
	return s
}

func (s *summary) Unit(unit string) *summary {
	s.unit = &unit
	return s
}

func (s *summary) Label(key, value string) *summary {
	newLabel := &dto.LabelPair{
		Name:  &key,
		Value: &value,
	}
	s.labels = append(s.labels, newLabel)
	return s
}

// Quantiles - please provide quantile values in ascending order!
func (s *summary) Quantiles(quantiles ...float64) *summary {
	s.quantiles = quantiles
	return s
}

func (s *summary) Observe(value float64) *summary {
	metric := s.storage.demand(s.name, s.help, s.unit, s.labels, dto.MetricType_SUMMARY, s.initMetrics)
	updateSummary(metric.Summary, value)
	return s
}

func (s *summary) initMetrics(metric *dto.Metric) {
	metric.Summary = &dto.Summary{
		SampleCount: new(uint64),
		SampleSum:   new(float64),
		Quantile:    makeQuantiles(s.quantiles),
	}
	// metric.Summary.quantileStore = newCKMS(s.quantiles)
}

func updateSummary(s *dto.Summary, value float64) {
	atomic.AddUint64(s.SampleCount, 1)
	AtomicAddFloat(s.SampleSum, value)
	// s.quantileStore.Insert(value)
	// for _, q := range s.Quantile {
	// 	*q.Value = s.quantileStore.Query(*q.Quantile)
	// }
}

func makeQuantiles(quantiles []float64) []*dto.Quantile {
	qvals := make([]float64, len(quantiles))
	quantileList := make([]*dto.Quantile, len(quantiles))
	for i, q := range quantiles {
		quantileList[i] = &dto.Quantile{
			Quantile: &q,
			Value:    &qvals[i],
		}
	}
	return quantileList
}
