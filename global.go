// @title           zeroprom metrics API
// @version         1.0
// @description     zerolog-inspired high performance metrics instrumentation API
// @termsOfService  http://swagger.io/terms/

// @contact.name   xakepp35
// @contact.url    https://github.com/xakepp35/zpm
// @contact.email  xakepp35@gmail.com

// @host      localhost:8080
// @BasePath  /metrics
package zpm

import (
	"io"

	"github.com/prometheus/common/expfmt"
)

// Export ➕
//
//	@Summary Exports metrics in a specified format.
//	@Description This function exports cumulative metrics, such as counters, in a desired format. Common use cases include exporting metrics for Prometheus scraping.
//	@Tags metrics
//	@Produce text/plain
//	@Param format query string true "The format to export the metrics in, e.g., 'text/plain'"
//	@Param options query string false "Additional encoding options"
//	@Usage Export metrics like processed jobs, HTTP request counts, etc.
//	@Misuse ❌ Using for fluctuating values (use Gauge for that).
//	@Misuse ❌ Attempting to decrease a counter (use Gauge instead).
//	@Pros ✅ Efficient for tracking incremental data.
//	@Pros ✅ Works well with `rate()` and `increase()` for Prometheus queries.
//	@Cons ⚠️ Cannot be decremented, reset is needed on restart.
//	@Tricks 🔍 Use `rate()` to track trends over time.
func Export(w io.Writer, expFormat expfmt.Format, opts ...expfmt.EncoderOption) error {
	return Srv.Export(w, expFormat, opts...)
}

// Counter ➕
//
//	@Summary Creates a new counter metric.
//	@Description This function creates a counter, a cumulative metric that only increases. Suitable for tracking things like requests served, errors, etc.
//	@Tags metrics
//	@Produce text/plain
//	@Param name query string true "Name of the counter metric"
//	@Usage Tracking events such as processed jobs or HTTP requests.
//	@Misuse ❌ Using for fluctuating values (use Gauge instead).
//	@Misuse ❌ Decreasing a counter (use Gauge for that).
//	@Pros ✅ Ideal for cumulative data.
//	@Pros ✅ Works seamlessly with `rate()` and `increase()` Prometheus functions.
//	@Cons ⚠️ Cannot be decremented.
//	@Tricks 🔍 Use `rate()` for trend analysis over time.
func Counter(name string) *counter {
	return Srv.Counter(name)
}

// Gauge ⚖️
//
//	@Summary Creates a new gauge metric.
//	@Description This function creates a gauge, a metric that can increase or decrease. Suitable for tracking values like temperature, memory usage, or concurrent requests.
//	@Tags metrics
//	@Produce text/plain
//	@Param name query string true "Name of the gauge metric"
//	@Usage Tracking instantaneous values like RAM usage, active users.
//	@Misuse ❌ Using for cumulative counts (use Counter instead).
//	@Pros ✅ Works well for instantaneous values.
//	@Pros ✅ Can both increase and decrease over time.
//	@Cons ⚠️ Quick fluctuations may make trend analysis challenging.
//	@Tricks 📊 Use `avg_over_time()` to smooth out variations and identify trends.
func Gauge(name string) *gauge {
	return Srv.Gauge(name)
}

// Histogram 📊
//
//	@Summary Creates a histogram metric for sampling observations.
//	@Description This function creates a histogram, which samples observations and categorizes them into configurable buckets. Useful for tracking durations or distributions.
//	@Tags metrics
//	@Produce text/plain
//	@Param name query string true "Name of the histogram metric"
//	@Usage Tracking response times, request sizes, or latencies.
//	@Misuse ❌ Too many buckets can lead to excessive memory usage.
//	@Pros ✅ Captures both the count and sum of observations.
//	@Pros ✅ Useful for estimating quantiles and understanding distributions.
//	@Cons ⚠️ Requires predefined bucket configuration. 
//	@Cons ⚠️ `histogram_quantile()` can produce misleading results with small samples.
//	@Tricks ⚡ Use the `le` label (`le="+Inf"`) to track total count. 🛠️ Use `rate()` on `_bucket` metrics for percentile estimations.
func Histogram(name string) *histogram {
	return Srv.Histogram(name)
}

// Summary 💡
//
//	@Summary Creates a summary metric for tracking dynamic quantiles.
//	@Description This function creates a summary, a metric that calculates quantiles and total observations on the client side. Ideal for latency tracking where dynamic quantiles are needed.
//	@Tags metrics
//	@Produce text/plain
//	@Param name query string true "Name of the summary metric"
//	@Usage Tracking request durations with precomputed quantiles, useful for latency analysis.
//	@Misuse ❌ Hard to aggregate across multiple instances (unlike Histograms).
//	@Pros ✅ Provides dynamic quantiles without predefined buckets.
//	@Cons ⚠️ Requires more memory than Histograms.
//	@Tricks 🎯 Use `quantile(0.95, rate(my_summary{quantile!=""}[5m]))` to estimate the 95th percentile of request durations.
func Summary(name string) *summary {
	return Srv.Summary(name)
}

// SortNames sets whether metric names should be ordered predictably during export.
//	@Summary Sets sorting behavior for metric names during export.
//	@Tags configuration
//	@Param sortNames query bool true "Whether to sort metric names predictably during export."
//	@Usage Enables deterministic ordering in exported metrics to facilitate easier comparisons.
func SortNames(sortNames bool) *Server {
	return Srv.SortNames(sortNames)
}

// String ➕
//
//	@Summary Exports metrics as a string in the specified format.
//	@Description This function exports the metrics as a string, supporting formats such as text/plain for Prometheus scraping or others as specified.
//	@Tags metrics
//	@Produce text/plain
//	@Param format query string true "The format to export the metrics in (e.g., 'text/plain')"
//	@Param options query string false "Additional encoding options"
//	@Usage Export metrics in a human-readable format or for Prometheus scraping.
//	@Misuse ❌ Using for non-cumulative metrics (use appropriate types).
//	@Tricks 📊 Use `rate()` in Prometheus queries for trend analysis.
func String(format expfmt.Format, opts ...expfmt.EncoderOption) (string, error) {
	return Srv.String(format, opts...)
}

// Bytes ➕
//
//	@Summary Exports metrics as raw bytes in the specified format.
//	@Description This function exports the metrics as raw bytes in a specified format, allowing for further processing or transport.
//	@Tags metrics
//	@Produce application/octet-stream
//	@Param format query string true "The format to export the metrics in (e.g., 'text/plain')"
//	@Param options query string false "Additional encoding options"
//	@Usage Export metrics for machine-to-machine communication or internal processing.
func Bytes(format expfmt.Format, opts ...expfmt.EncoderOption) (string, error) {
	return Srv.String(format, opts...)
}

// singletone
var Srv = NewServer()