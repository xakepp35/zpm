package zpm

import (
	dto "github.com/prometheus/client_model/go"
	"github.com/xakepp35/zpm/algo"
)

// Gauge client API interface
type gauge struct {
	name    string
	help    *string
	unit    *string
	labels  []*dto.LabelPair
	storage *storage
}

func (g *gauge) Help(help string) *gauge {
	g.help = &help
	return g
}

func (g *gauge) Unit(unit string) *gauge {
	g.unit = &unit
	return g
}

func (c *gauge) LabelPairs(labelPairs ...*LabelPair) *gauge {
	c.labels = append(c.labels, labelPairs...)
	return c
}

func (c *gauge) Label(key, value string) *gauge {
	return c.LabelPairs(&dto.LabelPair{
		Name:  &key,
		Value: &value,
	})
}

func (g *gauge) Set(value float64) *gauge {
	metricState := g.storage.demand(g.name, g.help, g.unit, g.labels, dto.MetricType_GAUGE, g.initMetric)
	algo.AtomicFloatStore(metricState.Dto.Gauge.Value, value)
	return g
}

func (g *gauge) Add(delta float64) *gauge {
	metricState := g.storage.demand(g.name, g.help, g.unit, g.labels, dto.MetricType_GAUGE, g.initMetric)
	algo.AtomicFloatAdd(metricState.Dto.Gauge.Value, delta)
	return g
}

func (g *gauge) Inc(delta int) *gauge {
	return g.Add(float64(delta))
}

func (g *gauge) Dec(delta int) *gauge {
	return g.Add(float64(-delta))
}

// initMetric creates a new gauge metric with labels
func (g *gauge) initMetric(metricState *state) {
	metricState.Dto.Gauge = &dto.Gauge{
		Value: new(float64),
	}
}
