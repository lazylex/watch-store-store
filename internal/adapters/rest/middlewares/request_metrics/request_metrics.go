package request_metrics

import (
	"github.com/lazylex/watch-store/store/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

type MiddlewareRequests struct {
	metrics *metrics.Metrics
}

// New конструктор прослойки для http-запросов
func New(metrics *metrics.Metrics) *MiddlewareRequests {
	return &MiddlewareRequests{metrics: metrics}
}

// RequestsInc увеличивает счетчик http-запросов к приложению. Добавляется метка path, содержащая путь запроса
func (m *MiddlewareRequests) RequestsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		m.metrics.HTTP.Requests.With(prometheus.Labels{metrics.PATH: r.RequestURI}).Inc()
		next.ServeHTTP(rw, r)
	})
}
