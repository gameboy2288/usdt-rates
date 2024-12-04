package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	rateRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_requests_total",
			Help: "Количество запросов к методу GetRates",
		},
		[]string{"status"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(rateRequests)
}

func RecordRequest(status string) {
	rateRequests.WithLabelValues(status).Inc()
}

func StartMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil) // Запускаем сервер метрик на 8080 порту
}
