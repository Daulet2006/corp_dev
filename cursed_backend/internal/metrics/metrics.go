package metrics

import (
	"cursed_backend/internal/logger"
	"expvar"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{Name: "http_requests_total", Help: "Total number of HTTP requests"},
		[]string{"method", "path", "status"},
	)
	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "Duration of HTTP requests", Buckets: prometheus.DefBuckets},
		[]string{"method", "path"},
	)
	errorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{Name: "http_errors_total", Help: "Total number of errors"},
		[]string{"type"},
	)
	DBQueryDuration = promauto.NewHistogramVec( // Exported (capital D)
		prometheus.HistogramOpts{Name: "db_query_duration_seconds", Help: "DB query duration"},
		[]string{"table"},
	)
	requestCount    = expvar.NewInt("requests_total")
	goroutinesCount = expvar.NewInt("goroutines_count")
)

func InitMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/debug/vars", expvar.Handler().ServeHTTP)
	// Health endpoint for alerting
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		goroutinesCount.Set(int64(runtime.NumGoroutine()))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status": "healthy", "goroutines": ` + strconv.Itoa(runtime.NumGoroutine()) + `}`))
		if err != nil {
			logger.Log.WithError(err).Error("Failed to write health check response")
			return
		}
	})
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		requestsTotal.WithLabelValues(method, path, status).Inc()
		requestDuration.WithLabelValues(method, path).Observe(duration)
		requestCount.Add(1)

		if c.Writer.Status() >= 400 {
			typ := "http"
			if strings.Contains(path, "/login") || strings.Contains(path, "/register") {
				typ = "auth"
			}
			errorsTotal.WithLabelValues(typ).Inc()
		}
	}
}
