package zpm

import (
	"math"
	"sync"
)

// CKMS algorithm for streaming quantile estimation
type ckms struct {
	sync.Mutex
	eps      float64
	quantile []float64
	stream   []pair
}

type pair struct {
	value float64
	g     int
	delta int
}

func newCKMS(quantiles []float64) *ckms {
	return &ckms{
		eps:      0.01, // accuracy trade-off
		quantile: quantiles,
		stream:   []pair{},
	}
}

func (c *ckms) Insert(value float64) {
	c.Lock()
	defer c.Unlock()
	c.stream = append(c.stream, pair{value, 1, 0})
	c.compress()
}

func (c *ckms) Query(q float64) float64 {
	c.Lock()
	defer c.Unlock()

	if len(c.stream) == 0 {
		return math.NaN()
	}

	rank := int(q * float64(len(c.stream)))
	if rank >= len(c.stream) {
		rank = len(c.stream) - 1
	}
	return c.stream[rank].value
}

func (c *ckms) compress() {
	// Simple compression strategy for CKMS
	if len(c.stream) < 2 {
		return
	}

	newStream := []pair{c.stream[0]}
	for i := 1; i < len(c.stream); i++ {
		last := &newStream[len(newStream)-1]
		curr := c.stream[i]
		if curr.value == last.value {
			last.g += curr.g
		} else {
			newStream = append(newStream, curr)
		}
	}
	c.stream = newStream
}
