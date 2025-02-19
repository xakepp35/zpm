package zpm

import (
	"sync/atomic"

	dto "github.com/prometheus/client_model/go"
	"github.com/xakepp35/zpm/algo"
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

func (s *summary) State() *state {
	return s.storage.demand(s.name, s.help, s.unit, s.labels, dto.MetricType_SUMMARY, s.initMetrics)
}

func (s *summary) Observe(value float64) *summary {
	metricState := s.State()
	atomic.AddUint64(metricState.Dto.Summary.SampleCount, 1)
	algo.AtomicFloatAdd(metricState.Dto.Summary.SampleSum, value)
	ckms := metricState.Data.(*algo.CKMSLockless)
	if ckms == nil {
		panic("CKMS is not initialized")
	}
	ckms.Insert(value)
	for _, q := range metricState.Dto.Summary.Quantile {
		*q.Value = ckms.Query(*q.Quantile)
	}
	return s
}

func (s *summary) initMetrics(metricState *state) {
	metricState.Dto.Summary = &dto.Summary{
		SampleCount: new(uint64),
		SampleSum:   new(float64),
		Quantile:    makeQuantiles(s.quantiles),
	}
	// CKMS для оценки квантилей
	metricState.Data = algo.NewCKMSLockless(s.quantiles...)
}

func updateSummary(s *dto.Summary, value float64) {

}

func makeQuantiles(quantiles []float64) []*dto.Quantile {
	qvals := make([]float64, len(quantiles))
	quantileList := make([]*dto.Quantile, len(quantiles))
	for i, q := range quantiles {
		quantileList[i] = &dto.Quantile{
			Quantile: &q,
			Value:    &qvals[i], // Изначально NaN
		}
	}
	return quantileList
}
