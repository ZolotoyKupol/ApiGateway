package metrics

import (
	"apigateway/internal/middleware"
	"net/http"
	"runtime"

	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var GoroutinesMetric = prometheus.NewGaugeFunc(
	prometheus.GaugeOpts{
		Name: "num_goroutines",
		Help: "Current number of goroutines",
	},
	func() float64 {
		return float64(runtime.NumGoroutine())
	},
)


var CacheSizeMetric = prometheus.NewGauge(
    prometheus.GaugeOpts{
        Name: "cache_size",
        Help: "Number of items in the cache",
    },
)

func InitMetrics(port string) {
	prometheus.MustRegister(GoroutinesMetric)
	prometheus.MustRegister(middleware.HttpStatusMetric)
	prometheus.MustRegister(CacheSizeMetric)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		slog.Debug("starting metrics server on port", "port", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			slog.Error("failed to start metrics server", "error", err)
		}
	}()

}


func UpdateCacheSizeMetric(size int) {
	CacheSizeMetric.Set(float64(size))
}