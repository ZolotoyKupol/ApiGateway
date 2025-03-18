package metrics

import (
	"apigateway/internal/models"
	"log/slog"
	"net/http"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	CacheSizeMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cache_size",
		Help: "Size of the cache",
	})

	CacheMemoryUsageMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cache_memory_usage",
		Help: "Memory usage of the cache",
	})
)

func CalculateCacheMemoryUsage(cache map[int]models.RoomDB) int64 {
	var totalSize int64

	totalSize += int64(unsafe.Sizeof(cache))

	for k, v := range cache {
		totalSize += int64(unsafe.Sizeof(k))
		totalSize += int64(unsafe.Sizeof(v))

		if v.Number != "" {
			totalSize += int64(len(v.Number))
		}
		if v.Floor != "" {
			totalSize += int64(len(v.Floor))
		}
		if v.Status != "" {
			totalSize += int64(len(v.Status))
		}
		if v.OccupiedBy != "" {
			totalSize += int64(len(v.OccupiedBy))
		}

		totalSize += int64(unsafe.Sizeof(v.CheckIn))
		totalSize += int64(unsafe.Sizeof(v.CheckOut))
	}

	return totalSize
}

var HttpStatusMetric = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_responses_total",
		Help: "Count of HTTP responses, labeled by status code and method",
	},
	[]string{"status", "method"},
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		statusCode := c.Writer.Status()
		method := c.Request.Method
		HttpStatusMetric.WithLabelValues(http.StatusText(statusCode), method).Inc()
		slog.Debug("HTTP request processed", "status", statusCode, "method", method)
	}
}

func Init(port string) {
	prometheus.MustRegister(HttpStatusMetric)
	prometheus.MustRegister(CacheSizeMetric)
	prometheus.MustRegister(CacheMemoryUsageMetric)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		slog.Debug("starting metrics server on port", "port", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			slog.Error("failed to start metrics server", "error", err)
		}
	}()
}

func UpdateCacheSizeMetric(cache map[int]models.RoomDB) {
	size := len(cache)
	memory := CalculateCacheMemoryUsage(cache)

	CacheSizeMetric.Set(float64(size))
	CacheMemoryUsageMetric.Set(float64(memory))
}
