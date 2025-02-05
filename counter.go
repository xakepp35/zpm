package zpm

import (
	dto "github.com/prometheus/client_model/go"
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

func (c *counter) Label(key, value string) *counter {
	newLabel := &dto.LabelPair{
		Name:  &key,
		Value: &value,
	}
	c.labels = append(c.labels, newLabel)
	return c
}

func (c *counter) Add(delta float64) *counter {
	metric := c.storage.demand(c.name, c.help, c.unit, c.labels, dto.MetricType_COUNTER, c.initMetric)
	AtomicAddFloat(metric.Counter.Value, delta)
	return c
}

func (c *counter) Inc(delta int) *counter {
	return c.Add(float64(delta))
}

// newMetric creates a new counter metric with labels
func (c *counter) initMetric(metric *dto.Metric) {
	value := float64(0)
	metric.Counter = &dto.Counter{
		Value: &value,
	}
}
