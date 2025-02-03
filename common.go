package zpm

import (
	"math"
	"sync/atomic"
	"time"
	"unsafe"

	dto "github.com/prometheus/client_model/go"
)

func makeKey(name string, labels []*dto.LabelPair) string {
	key := name
	for _, lbl := range labels {
		key += "|" + *lbl.Value
	}
	return key
}

// newMetric creates a new counter metric with labels
func newMetricCounter(timestampMs int64, labels []*dto.LabelPair) *dto.Metric {
	value := float64(0)
	return &dto.Metric{
		Label: labels,
		Counter: &dto.Counter{
			Value: &value,
		},
		TimestampMs: &timestampMs,
	}
}

func TimestampMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// AtomicAddFloat is atomic float addition using CAS
func AtomicAddFloat(x *float64, delta float64) {
	addr := (*uint64)(unsafe.Pointer(x))
	for {
		oldBits := atomic.LoadUint64(addr)
		oldVal := math.Float64frombits(oldBits)
		newVal := oldVal + delta
		newBits := math.Float64bits(newVal)
		if atomic.CompareAndSwapUint64(addr, oldBits, newBits) {
			return
		}
	}
}
