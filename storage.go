package zpm

import (
	"fmt"
	"sort"
	"sync"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

type storage struct {
	mu       sync.RWMutex
	metrics  map[string]*dto.Metric
	families map[string]*dto.MetricFamily
	names    []string
}

func NewStorage() *storage {
	return &storage{
		metrics:  make(map[string]*dto.Metric),
		families: make(map[string]*dto.MetricFamily),
	}
}

func (s *storage) demand(name string, help, unit *string, labels []*dto.LabelPair, metricType dto.MetricType, initMetric func(metric *dto.Metric)) *dto.Metric {
	key := makeKey(name, labels)
	metric := s.get(key)
	if metric != nil {
		return metric
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	metric = s.metrics[key]
	if metric != nil {
		return metric
	}
	fam := s.families[name]
	if fam == nil {
		fam = &dto.MetricFamily{
			Name: &name,
			Help: help,
			Type: metricType.Enum(),
			Unit: unit,
		}
		s.families[name] = fam
		s.names = append(s.names, name)
	}
	timestampMs := TimestampMs()
	metric = &dto.Metric{
		Label:       labels,
		TimestampMs: &timestampMs,
	}
	initMetric(metric)
	metric.Label = labels
	fam.Metric = append(fam.Metric, metric)
	s.metrics[key] = metric
	return metric
}

func (s *storage) get(key string) *dto.Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.metrics[key]
}

func (s *storage) Encode(encoder expfmt.Encoder, sortNames bool) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if sortNames {
		sort.Strings(s.names)
	}
	for _, name := range s.names {
		metricFamily := s.families[name]
		if err := encoder.Encode(metricFamily); err != nil {
			return fmt.Errorf("expfmt.Encode(%s): %v", *metricFamily.Name, err)
		}
	}
	return nil
}
