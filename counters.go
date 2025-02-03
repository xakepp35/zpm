package zpm

import (
	"fmt"
	"sort"
	"sync"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

// counters manages named counters and histograms
type counters struct {
	mu        sync.RWMutex
	counters  map[string]*dto.Metric
	families  map[string]*dto.MetricFamily
	names     []string
	sortNames bool
}

// NewCollectors initializes a new Collectors instance with a specified format
func NewCounters() *counters {
	return &counters{
		counters: make(map[string]*dto.Metric),
		families: make(map[string]*dto.MetricFamily),
	}
}

func (c *counters) Counter(name string) *counter {
	return &counter{
		name: name,
		m:    c,
	}
}

// getMetric retrieves the metric if it exists
func (c *counters) getCounter(key string) *dto.Metric {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.counters[key]
}

// demandMetric ensures the metric exists, creating it if necessary
func (c *counters) demand(m *counter) *dto.Metric {
	key := makeKey(m.name, m.labels)
	metric := c.getCounter(key)
	if metric != nil {
		return metric
	}
	// first call will follow slow path
	c.mu.Lock()
	defer c.mu.Unlock()
	// double-check in case another goroutine created it
	metric = c.counters[key]
	if metric != nil {
		return metric
	}
	// create new
	metric = newMetricCounter(TimestampMs(), m.labels)
	// check family registration
	fam := c.families[m.name]
	// first metricFamily creation
	if fam == nil {
		fam = &dto.MetricFamily{
			Name: &m.name,
			Help: m.help,
			Type: dto.MetricType_COUNTER.Enum(),
			Unit: m.unit,
		}
		c.families[m.name] = fam
		c.names = append(c.names, m.name)
	}
	fam.Metric = append(fam.Metric, metric)
	// named+labelled metric registration
	c.counters[key] = metric
	return metric
}

// ExportMetrics outputs metrics in Prometheus format using expfmt; pass expfmt.TypeTextPlain as second
func (c *counters) Encode(encoder expfmt.Encoder) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.sortNames {
		sort.Strings(c.names)
	}
	for _, name := range c.names {
		metricFamily := c.families[name]
		err := encoder.Encode(metricFamily)
		if err != nil {
			return fmt.Errorf("counters.encode(%s): %v", *metricFamily.Name, err)
		}
	}
	return nil
}
