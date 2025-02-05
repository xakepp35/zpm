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

func TimestampMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// atomically stores float64 newVal into x
func AtomicSetFloat(x *float64, newVal float64) {
	addr := (*uint64)(unsafe.Pointer(x))
	newBits := math.Float64bits(newVal)
	atomic.StoreUint64(addr, newBits)
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
