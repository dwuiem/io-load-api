package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"log/slog"
	"net/http"
	"test-workmate/internal/config"
)

var (
	HttpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Time duration of HTTP Request",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status_code"},
	)

	TaskProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tasks_processed_total",
			Help: "Total number of tasks processed",
		},
		[]string{"status"},
	)

	ActiveTasks = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_tasks",
			Help: "Total number of active tasks",
		},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(TaskProcessed)
	prometheus.MustRegister(ActiveTasks)
	prometheus.MustRegister(HttpDuration)
}

func StartMetricsServer(cfg *config.Config) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.PrometheusPort), nil)
		if err != nil {
			log.Fatal("Metrics server failed", slog.String("error", err.Error()))
		}
	}()
}
