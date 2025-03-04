package zpm

import (
	"sync/atomic"

	dto "github.com/prometheus/client_model/go"
	"github.com/xakepp35/zpm/algo"
)

// Histogram client API
type histogram struct {
	name    string
	help    *string
	unit    *string
	labels  []*dto.LabelPair
	buckets []float64
	storage *storage
}

func (h *histogram) Help(help string) *histogram {
	h.help = &help
	return h
}

func (h *histogram) Unit(unit string) *histogram {
	h.unit = &unit
	return h
}

func (c *histogram) LabelPairs(labelPairs ...*LabelPair) *histogram {
	c.labels = append(c.labels, labelPairs...)
	return c
}

func (c *histogram) Label(key, value string) *histogram {
	return c.LabelPairs(&dto.LabelPair{
		Name:  &key,
		Value: &value,
	})
}

// Buckets - please provide sorted bucket values in ascending order!
func (h *histogram) Buckets(buckets ...float64) *histogram {
	h.buckets = buckets
	return h
}

func (h *histogram) Observe(value float64) *histogram {
	metricState := h.storage.demand(h.name, h.help, h.unit, h.labels, dto.MetricType_HISTOGRAM, h.initMetric)
	updateHistogram(metricState.Dto.Histogram, value)
	return h
}

func (h *histogram) initMetric(metricState *state) {
	metricState.Dto.Histogram = &dto.Histogram{
		SampleCount: new(uint64),
		SampleSum:   new(float64),
		Bucket:      makeBuckets(h.buckets),
	}
}

func updateHistogram(h *dto.Histogram, value float64) {
	atomic.AddUint64(h.SampleCount, 1)
	algo.AtomicFloatAdd(h.SampleSum, value)
	for _, b := range h.Bucket {
		if value <= *b.UpperBound {
			atomic.AddUint64(b.CumulativeCount, 1)
		}
	}
}

func makeBuckets(bounds []float64) []*dto.Bucket {
	counts := make([]uint64, len(bounds))
	buckets := make([]*dto.Bucket, len(bounds))
	for i := range bounds {
		buckets[i] = &dto.Bucket{
			UpperBound:      &bounds[i],
			CumulativeCount: &counts[i],
		}
	}
	return buckets
}
