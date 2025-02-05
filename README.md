# zeroprom - zerolog-inspired High-Performance Prometheus Metrics Wrapper

`zpm` is a minimalistic and efficient Prometheus metrics wrapper for Go, designed for high-performance applications. It simplifies counter and histogram metric handling while ensuring concurrency safety.

## Features  

✅ **Simplicity** – Easy-to-use API with minimal setup  
🚀 **Efficiency** – Optimized for high-performance applications  
🔒 **Thread-Safe** – Designed for concurrent access in multithreaded applications  
🔖 **Labeled** – Supports labeled counters and histograms for better observability  
📡 **Compatible** – Exposes metrics in a Prometheus-friendly format  
🛠 **Customizable** – Supports different output formats using `expfmt`  


## Installation

```sh
go get github.com/xakepp35/zpm
```

## Examples

```go
// measure and perform your call
startedAt := time.Now()
res, err := PerformMyRequest()
latencyMs := time.Since(startedAt).Seconds()/1000

// zerolog example, for visual comparison:
log.Info().
    Err(err).
    Any("res", res).
    Str("func", zpm.RuntimeFunctionName(0)).
    Msg("request")

// counter example:
zpm.Counter("http_requests_total").
    Help("http requests counter").
    Label("method", r.Method).
    Label("path", r.URL.Path).
    Inc(1)

// gauge example:
zpm.Gauge("http_requests_gauge").
    Help("http requests latency gauge").
    Label("method", r.Method).
    Label("path", r.URL.Path).
    Set(latencyMs)

// histogrtam example:
zpm.Histogram("http_duration_milliseconds").
    Help("http requests duration histogram").
    Buckets(1, 10, 100, 1000).
    Label("method", r.Method).
    Label("path", r.URL.Path).
    Observe(latencyMs)

// summary example:
zpm.Summary("http_duration_summary_milliseconds")
    Help("http requests duration summary").
    Quantiles(0, 0.1, 0.5, 0.9, 1).
    Label("method", r.Method).
    Label("path", r.URL.Path).
    Observe(latencyMs)
```

## License

This project is licensed under the MIT License.