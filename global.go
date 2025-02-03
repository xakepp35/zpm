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

// func Histograms(name string) *counter {
// 	return zpc.Histogram(name)
// }
