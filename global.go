package zpm

import (
	"io"

	"github.com/prometheus/common/expfmt"
)

var gServer = NewServer()

func SortNames(sortNames bool) *Server {
	return gServer.SortNames(sortNames)
}

func String(opts ...expfmt.EncoderOption) (string, error) {
	return gServer.String(opts...)
}

func Export(w io.Writer, expFormat expfmt.Format, opts ...expfmt.EncoderOption) error {
	return gServer.Export(w, expFormat, opts...)
}

func Counter(name string) *counter {
	return gServer.Counter(name)
}

func Gauge(name string) *gauge {
	return gServer.Gauge(name)
}

func Histogram(name string) *histogram {
	return gServer.Histogram(name)
}

func Summary(name string) *summary {
	return gServer.Summary(name)
}
