package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type Metrics struct {
	opts                       prometheus.Opts
	registry                   *prometheus.Registry
	httpRequestTotal           *prometheus.CounterVec
	httpRequestDurationSeconds *prometheus.HistogramVec
}

func New(svc string, ops ...Option) *Metrics {

	m := &Metrics{
		opts: prometheus.Opts{
			ConstLabels: prometheus.Labels{
				"svc": svc,
			},
		},
		registry: prometheus.NewRegistry(),
	}
	for _, op := range ops {
		op(m)
	}

	m.httpRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        "http_request_total",
		Help:        "HTTP 请求计数",
		ConstLabels: m.opts.ConstLabels,
	}, []string{"method", "path"})

	m.httpRequestDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "http_request_duration_seconds",
		Help:        "HTTP 请求响应时长（秒）",
		ConstLabels: m.opts.ConstLabels,
	}, []string{"method", "path"})

	m.registry.MustRegister(
		m.httpRequestTotal,
		m.httpRequestDurationSeconds,
	)
	return m
}

func (m *Metrics) IncRequestTotal(method, path string) {
	m.httpRequestTotal.WithLabelValues(method, path).Inc()
}

func (m *Metrics) ObserveRequestDuration(method, path string, d time.Duration) {
	m.httpRequestDurationSeconds.WithLabelValues(method, path).Observe(d.Seconds())
}

// Gin 专用中间件
func (m *Metrics) Gin(g *gin.Context) {
	t := time.Now()
	g.Next()
	d := time.Since(t)
	method := g.Request.Method
	path := g.Request.URL.Path
	m.IncRequestTotal(method, path)
	m.ObserveRequestDuration(method, path, d)
}

// Exporter prometheus 指标收集接口
// gin 示例 r.GET("/metrics", gin.WrapH(m.Exporter()))
func (m *Metrics) Exporter() http.Handler {
	return promhttp.InstrumentMetricHandler(
		m.registry,
		promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}),
	)
}
