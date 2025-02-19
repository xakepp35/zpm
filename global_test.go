package zpm_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xakepp35/zpm"
)

func TestHistograms(t *testing.T) {
	const numIter = 100
	var wg sync.WaitGroup
	wg.Add(2 * numIter)
	for i := 0; i < numIter; i++ {
		go func() {
			defer wg.Done()
			zpm.Histogram("hist_1").
				Help("hist1 help").
				Unit("hist1 unit").
				Buckets(0.1, 1, 10).
				Label("l1", "v1").
				Label("l2", "v2").
				Observe(0.1 * float64(i))
		}()
		go func() {
			defer wg.Done()
			zpm.Histogram("hist_1").
				Help("hist1 help").
				Unit("hist1 unit").
				Buckets(0.1, 1, 10).
				Label("l1", "v2").
				Label("l2", "v2").
				Observe(0.2 * float64(i))
		}()
	}
	wg.Wait()
	res, err := zpm.String()
	fmt.Printf("%s", res)
	assert.NoError(t, err)
	assert.NotEqual(t, "", res)
}

func TestCounters(t *testing.T) {
	const numIter = 1000
	var wg sync.WaitGroup
	wg.Add(4 * numIter)
	zpm.SortNames(true)
	for i := 0; i < numIter; i++ {
		go func() {
			defer wg.Done()
			zpm.Counter("ctr_1").
				Help("ctr1 help").
				Unit("ctr1 unit").
				Label("l1", "v1").
				Label("l2", "v2").
				Inc(1)
		}()
		go func() {
			defer wg.Done()
			zpm.Counter("ctr_2").
				Help("ctr2 help").
				Label("l1", "v1").
				Label("l2", "v2").
				Inc(1)
		}()
		go func() {
			defer wg.Done()
			zpm.Counter("ctr_2").
				Help("ctr3 help").
				Label("l1", "v2").
				Label("l2", "v2").
				Inc(1)
		}()
		go func() {
			defer wg.Done()
			zpm.Counter("ctr_2").
				Help("ctr2 help").
				Label("l1", "v1").
				Label("l2", "v2").
				Inc(1)
		}()
	}
	wg.Wait()
	res, err := zpm.String()
	fmt.Printf("%s", res)
	assert.NoError(t, err)
	assert.NotEqual(t, "", res)
}
