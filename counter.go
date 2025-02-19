package zpm

import (
	dto "github.com/prometheus/client_model/go"
	"github.com/xakepp35/zpm/algo"
)

// Counter client API interface
type counter struct {
	name    string
	help    *string
	unit    *string
	labels  []*dto.LabelPair
	storage *storage
}

func (c *counter) Help(help string) *counter {
	c.help = &help
	return c
}

func (c *counter) Unit(unit string) *counter {
	c.unit = &unit
	return c
}

func (c *counter) LabelPairs(labelPairs ...*LabelPair) *counter {
	c.labels = append(c.labels, labelPairs...)
	return c
}

func (c *counter) Label(key, value string) *counter {
	return c.LabelPairs(&dto.LabelPair{
		Name:  &key,
		Value: &value,
	})
}

// Please, be careful: counter should be everincreasing value!
func (c *counter) Set(value float64) *counter {
	metricState := c.storage.demand(c.name, c.help, c.unit, c.labels, dto.MetricType_COUNTER, c.initMetric)
	algo.AtomicFloatStore(metricState.Dto.Counter.Value, value)
	return c
}

func (c *counter) Add(delta float64) *counter {
	metricState := c.storage.demand(c.name, c.help, c.unit, c.labels, dto.MetricType_COUNTER, c.initMetric)
	algo.AtomicFloatAdd(metricState.Dto.Counter.Value, delta)
	return c
}

func (c *counter) Inc(delta int) *counter {
	return c.Add(float64(delta))
}

// newMetric creates a new counter metric with labels
func (c *counter) initMetric(metricState *state) {
	value := float64(0)
	metricState.Dto.Counter = &dto.Counter{
		Value: &value,
	}
}
