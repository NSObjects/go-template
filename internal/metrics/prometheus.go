package metrics

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusMetrics Prometheus指标收集器
type PrometheusMetrics struct {
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	httpRequestSize     *prometheus.HistogramVec
	httpResponseSize    *prometheus.HistogramVec
	activeConnections   prometheus.Gauge
	databaseConnections *prometheus.GaugeVec
	cacheHits           *prometheus.CounterVec
	cacheMisses         *prometheus.CounterVec
}

// NewPrometheusMetrics 创建Prometheus指标收集器
func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status_code"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status_code"},
		),
		httpRequestSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_size_bytes",
				Help:    "HTTP request size in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "path"},
		),
		httpResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "HTTP response size in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "path", "status_code"},
		),
		activeConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_connections",
				Help: "Number of active connections",
			},
		),
		databaseConnections: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "database_connections",
				Help: "Number of database connections",
			},
			[]string{"type", "status"},
		),
		cacheHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_hits_total",
				Help: "Total number of cache hits",
			},
			[]string{"cache_type", "operation"},
		),
		cacheMisses: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_misses_total",
				Help: "Total number of cache misses",
			},
			[]string{"cache_type", "operation"},
		),
	}
}

// HTTPMiddleware HTTP指标中间件
func (m *PrometheusMetrics) HTTPMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// 记录请求大小
			requestSize := c.Request().ContentLength
			if requestSize > 0 {
				m.httpRequestSize.WithLabelValues(
					c.Request().Method,
					c.Path(),
				).Observe(float64(requestSize))
			}

			// 处理请求
			err := next(c)

			// 记录响应
			duration := time.Since(start).Seconds()
			statusCode := strconv.Itoa(c.Response().Status)

			m.httpRequestsTotal.WithLabelValues(
				c.Request().Method,
				c.Path(),
				statusCode,
			).Inc()

			m.httpRequestDuration.WithLabelValues(
				c.Request().Method,
				c.Path(),
				statusCode,
			).Observe(duration)

			// 记录响应大小
			responseSize := c.Response().Size
			if responseSize > 0 {
				m.httpResponseSize.WithLabelValues(
					c.Request().Method,
					c.Path(),
					statusCode,
				).Observe(float64(responseSize))
			}

			return err
		}
	}
}

// RecordCacheHit 记录缓存命中
func (m *PrometheusMetrics) RecordCacheHit(cacheType, operation string) {
	m.cacheHits.WithLabelValues(cacheType, operation).Inc()
}

// RecordCacheMiss 记录缓存未命中
func (m *PrometheusMetrics) RecordCacheMiss(cacheType, operation string) {
	m.cacheMisses.WithLabelValues(cacheType, operation).Inc()
}

// SetActiveConnections 设置活跃连接数
func (m *PrometheusMetrics) SetActiveConnections(count int) {
	m.activeConnections.Set(float64(count))
}

// SetDatabaseConnections 设置数据库连接数
func (m *PrometheusMetrics) SetDatabaseConnections(dbType, status string, count int) {
	m.databaseConnections.WithLabelValues(dbType, status).Set(float64(count))
}
