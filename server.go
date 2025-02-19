package zpm

import (
	"bytes"
	"fmt"
	"io"

	"github.com/prometheus/common/expfmt"
)

var (
	FmtTextPlain = expfmt.NewFormat(expfmt.TypeTextPlain)
)

type ServerConfig struct {
	SortNames bool `json:"sort_names"`
}

type Server struct {
	counters   *storage
	gauges     *storage
	histograms *storage
	summaries  *storage

	cfg *ServerConfig
}

func (s *Server) OptSortNames(sortNames bool) *Server {
	s.cfg.SortNames = sortNames
	return s
}

func (s *Server) Counter(name string) *counter {
	return &counter{
		name:    name,
		storage: s.counters,
	}
}

func (s *Server) Gauge(name string) *gauge {
	return &gauge{
		name:    name,
		storage: s.gauges,
	}
}

func (s *Server) Histogram(name string) *histogram {
	return &histogram{
		name:    name,
		storage: s.histograms,
	}
}

func (s *Server) Summary(name string) *summary {
	return &summary{
		name:    name,
		storage: s.summaries,
	}
}

func NewServer() *Server {
	return &Server{
		counters:   NewStorage(),
		gauges:     NewStorage(),
		histograms: NewStorage(),
		summaries:  NewStorage(),
		cfg: &ServerConfig{
			SortNames: false,
		},
	}
}

func (s *Server) SortNames(sortNames bool) *Server {
	s.cfg.SortNames = sortNames
	return s
}

func (s *Server) Export(w io.Writer, expFormat expfmt.Format, opts ...expfmt.EncoderOption) error {
	encoder := expfmt.NewEncoder(w, expFormat, opts...)
	if err := s.counters.Encode(encoder, s.cfg.SortNames); err != nil {
		return fmt.Errorf("counters.Encode(): %w", err)
	}
	if err := s.histograms.Encode(encoder, s.cfg.SortNames); err != nil {
		return fmt.Errorf("counters.Encode(): %w", err)
	}
	return nil
}

// String renders the Prometheus texp/plain formatted metrics to string
func (c *Server) String(format expfmt.Format, opts ...expfmt.EncoderOption) (string, error) {
	var buf bytes.Buffer
	if err := c.Export(&buf, format, opts...); err != nil {
		return "", err
	}
	res := buf.String()
	return res, nil
}

// String renders the Prometheus texp/plain formatted metrics to string
func (c *Server) Bytes(format expfmt.Format, opts ...expfmt.EncoderOption) ([]byte, error) {
	var buf bytes.Buffer
	if err := c.Export(&buf, format, opts...); err != nil {
		return nil, err
	}
	res := buf.Bytes()
	return res, nil
}
