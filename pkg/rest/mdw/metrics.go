package mdw

import (
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

var (
	// Latency (response time)
	httpResponseLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "Latency",
			Help:    "Request latency distribution (successful and error responses)",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 13),
		},
		[]string{"path", "status_class"},
	)

	// Traffic (request count)
	httpTrafficTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "Traffic",
			Help: "Total number of HTTP requests (traffic)",
		},
		[]string{"path", "status_class"},
	)

	// Errors
	httpErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "Errors",
			Help: "Total number of HTTP error responses (4xx and 5xx)",
		},
		[]string{"path", "status_code"},
	)

	// Saturation (goroutines + heap)
	appGoroutines = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "Saturation_goroutines",
			Help: "Number of running goroutines (application saturation)",
		},
		func() float64 {
			return float64(runtime.NumGoroutine())
		},
	)

	appHeapAlloc = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "Saturation_heap",
			Help: "Allocated heap memory in bytes",
		},
		func() float64 {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			return float64(m.HeapAlloc)
		},
	)

	//Request per seconds
	requestPerSecondBuckets = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "RequestPerSecond",
			Help: "Histogram of requests per second over response time buckets",
			Buckets: []float64{
				0.001, 0.005, 0.01, 0.05,
				0.1, 0.2, 0.3, 0.5,
				1.0, 2.0, 3.0, 4.0, 5.0,
			},
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(
		httpResponseLatency,
		httpTrafficTotal,
		httpErrorsTotal,
		appGoroutines,
		appHeapAlloc,
		requestPerSecondBuckets,
	)
}

func MetricsMiddleware(next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next(rr, r)

		duration := time.Since(start).Seconds()
		statusClass := getStatusClass(rr.statusCode)

		// Latency
		httpResponseLatency.WithLabelValues(r.URL.Path, statusClass).Observe(duration)

		// Traffic
		httpTrafficTotal.WithLabelValues(r.URL.Path, statusClass).Inc()

		// Request per second
		requestPerSecondBuckets.WithLabelValues(r.URL.Path).Observe(duration)

		// Errors
		if rr.statusCode >= 400 {
			httpErrorsTotal.WithLabelValues(r.URL.Path, http.StatusText(rr.statusCode)).Inc()
		}
	}
}

func getStatusClass(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "2xx"
	case code >= 300 && code < 400:
		return "3xx"
	case code >= 400 && code < 500:
		return "4xx"
	case code >= 500:
		return "5xx"
	default:
		return "unknown"
	}
}
