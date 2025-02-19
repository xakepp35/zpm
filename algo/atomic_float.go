package algo

import (
	"math"
	"sync/atomic"
	"unsafe"
)

// atomically stores float64 newVal into x
func AtomicFloatStore(x *float64, newVal float64) {
	addr := (*uint64)(unsafe.Pointer(x))
	bits := math.Float64bits(newVal)
	atomic.StoreUint64(addr, bits)
}

func AtomicFloatLoad(x *float64) float64 {
	addr := (*uint64)(unsafe.Pointer(x))
	bits := atomic.LoadUint64(addr)
	return math.Float64frombits(bits)
}

// AtomicAddFloat is atomic float addition using CAS
func AtomicFloatAdd(x *float64, delta float64) {
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
