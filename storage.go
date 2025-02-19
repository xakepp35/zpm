package zpm

import (
	"fmt"
	"sort"
	"sync"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/xakepp35/zpm/algo"
)

type storage struct {
	mu       sync.RWMutex
	metrics  map[string]*state
	families map[string]*dto.MetricFamily
	names    []string
}

func NewStorage() *storage {
	return &storage{
		metrics:  make(map[string]*state),
		families: make(map[string]*dto.MetricFamily),
	}
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

func (s *storage) demand(name string, help, unit *string, labels []*dto.LabelPair, metricType dto.MetricType, initMetric StateInitFunc) *state {
	key := makeKey(name, labels)
	metricState := s.get(key)
	if metricState != nil {
		return metricState
	}
	// first time: init once
	s.mu.Lock()
	defer s.mu.Unlock()
	metricState = s.metrics[key]
	if metricState != nil {
		return metricState
	}
	fam := s.families[name]
	if fam == nil {
		fam = s.registerName(name, help, unit, metricType)
	}
	timestampMs := algo.TimestampMs()
	metricState = newState(timestampMs, labels...)
	initMetric(metricState)
	metricState.Dto.Label = labels
	fam.Metric = append(fam.Metric, metricState.Dto)
	s.metrics[key] = metricState
	return metricState
}

func (s *storage) registerName(name string, help, unit *string, metricType dto.MetricType) *dto.MetricFamily {
	res := &dto.MetricFamily{
		Name: &name,
		Help: help,
		Type: metricType.Enum(),
		Unit: unit,
	}
	s.families[name] = res
	s.names = append(s.names, name)
	return res
}

const labelsSeparator = "\u001d"

func makeKey(name string, labels []*dto.LabelPair) string {
	key := name
	for _, lbl := range labels {
		key += labelsSeparator + *lbl.Value
	}
	return key
}

func (s *storage) get(key string) *state {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.metrics[key]
}
