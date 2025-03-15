package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
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

var HttpStatusMetric = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_responses_total",
		Help: "Count of HTTP responses, labeled by status code and method",
	},
	[]string{"status", "method"},
)