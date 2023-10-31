package request_metrics

import (
	"github.com/lazylex/watch-store/store/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strings"
)

type MiddlewareRequests struct {
	metrics *metrics.Metrics
}

// New конструктор прослойки для http-запросов
func New(metrics *metrics.Metrics) *MiddlewareRequests {
	return &MiddlewareRequests{metrics: metrics}
}

// RequestsInc увеличивает счетчик http-запросов к приложению. Добавляется метка path, содержащая путь запроса.
// Метка path для GET-запросов сохраняется без идентификатора
func (m *MiddlewareRequests) RequestsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			m.metrics.HTTP.Requests.With(prometheus.Labels{
				metrics.PATH: string([]rune(r.RequestURI)[0 : strings.LastIndex(r.RequestURI, "/")+1]),
			}).Inc()
		} else {
			m.metrics.HTTP.Requests.With(prometheus.Labels{metrics.PATH: r.RequestURI}).Inc()
		}
		next.ServeHTTP(rw, r)
	})
}
