package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	// Технические метрики
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_response_time_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: []float64{0.1, 0.3, 0.5, 1, 2, 5},
	}, []string{"method", "path"})

	// Бизнесовые метрики
	pvzCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pvz_created_total",
		Help: "Total number of PVZ created",
	})

	receptionsCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "receptions_created_total",
		Help: "Total number of receptions created",
	})

	productsAdded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_added_total",
		Help: "Total number of products added",
	})
)

func InitPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":9000", nil)
	}()
}

func RecordHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	httpRequestsTotal.WithLabelValues(method, path, http.StatusText(statusCode)).Inc()
	httpDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

func RecordPVZCreated() {
	pvzCreated.Inc()
}

func RecordReceptionCreated() {
	receptionsCreated.Inc()
}

func RecordProductAdded() {
	productsAdded.Inc()
}
