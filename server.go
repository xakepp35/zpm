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

type Server struct {
	counters *counters
}

func (s *Server) Counter(name string) *counter {
	return s.counters.Counter(name)
}

func NewServer() *Server {
	return &Server{
		counters: NewCounters(),
	}
}

func (s *Server) SortNames(sortNames bool) *Server {
	s.counters.sortNames = sortNames
	return s
}

func (s *Server) Export(w io.Writer, expFormat expfmt.Format, opts ...expfmt.EncoderOption) error {
	encoder := expfmt.NewEncoder(w, expFormat, opts...)
	if err := s.counters.Encode(encoder); err != nil {
		return fmt.Errorf("counters.Encode(): %w", err)
	}
	return nil
}

// String renders the Prometheus texp/plain formatted metrics to string
func (c *Server) String(opts ...expfmt.EncoderOption) (string, error) {
	var buf bytes.Buffer
	if err := c.Export(&buf, FmtTextPlain, opts...); err != nil {
		return "", err
	}
	res := buf.String()
	return res, nil
}

// String renders the Prometheus texp/plain formatted metrics to string
func (c *Server) Bytes(opts ...expfmt.EncoderOption) ([]byte, error) {
	var buf bytes.Buffer
	if err := c.Export(&buf, FmtTextPlain, opts...); err != nil {
		return nil, err
	}
	res := buf.Bytes()
	return res, nil
}
