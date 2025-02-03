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

## Example

```go
// remember when you have started your call
startedAt := time.Now()

// perform http request here
res, err := PerformMyRequest()

// log with zerolog:
log.Info().
    Any("res", res).
    Err(err).
    Str("func", zpm.RuntimeFunctionName(0)).
    Msg("request")

// measure requests count:
zpm.Counter("http_requests_total").
    Help("http requests counter")
    Label("method", r.Method).
    Label("path", r.URL.Path).
    Inc(1)

// measure histogrtam latencies:
zpm.Histogram("http_duration_milliseconds").
    Help("http requests duration")
    Buckets(1, 10, 100, 1000).
    Label("method", r.Method).
    Label("path", r.URL.Path).
    Observe(time.Since(startedAt).Seconds()/1000)
```

## License

This project is licensed under the MIT License.