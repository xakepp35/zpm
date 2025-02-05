package zpm

import (
	dto "github.com/prometheus/client_model/go"
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

func (g *gauge) Label(key, value string) *gauge {
	newLabel := &dto.LabelPair{
		Name:  &key,
		Value: &value,
	}
	g.labels = append(g.labels, newLabel)
	return g
}

func (g *gauge) Set(value float64) *gauge {
	metric := g.storage.demand(g.name, g.help, g.unit, g.labels, dto.MetricType_GAUGE, g.initMetric)
	AtomicSetFloat(metric.Gauge.Value, value)
	return g
}

func (g *gauge) Add(delta float64) *gauge {
	metric := g.storage.demand(g.name, g.help, g.unit, g.labels, dto.MetricType_GAUGE, g.initMetric)
	AtomicAddFloat(metric.Gauge.Value, delta)
	return g
}

func (g *gauge) Inc(delta int) *gauge {
	return g.Add(float64(delta))
}

func (g *gauge) Dec(delta int) *gauge {
	return g.Add(float64(-delta))
}

// initMetric creates a new gauge metric with labels
func (g *gauge) initMetric(metric *dto.Metric) {
	metric.Gauge = &dto.Gauge{
		Value: new(float64),
	}
}
