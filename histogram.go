package zpm

import (
	"sync/atomic"

	dto "github.com/prometheus/client_model/go"
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

func (h *histogram) Label(key, value string) *histogram {
	newLabel := &dto.LabelPair{
		Name:  &key,
		Value: &value,
	}
	h.labels = append(h.labels, newLabel)
	return h
}

// Buckets - please provide sorted bucket values in ascending order!
func (h *histogram) Buckets(buckets ...float64) *histogram {
	h.buckets = buckets
	return h
}

func (h *histogram) Observe(value float64) *histogram {
	metric := h.storage.demand(h.name, h.help, h.unit, h.labels, dto.MetricType_HISTOGRAM, h.initMetrics)
	updateHistogram(metric.Histogram, value)
	return h
}

func (h *histogram) initMetrics(metric *dto.Metric) {
	metric.Histogram = &dto.Histogram{
		SampleCount: new(uint64),
		SampleSum:   new(float64),
		Bucket:      makeBuckets(h.buckets),
	}
}

func updateHistogram(h *dto.Histogram, value float64) {
	atomic.AddUint64(h.SampleCount, 1)
	AtomicAddFloat(h.SampleSum, value)
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
