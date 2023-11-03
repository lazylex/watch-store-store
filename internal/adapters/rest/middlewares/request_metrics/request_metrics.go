package request_metrics

import (
	"github.com/lazylex/watch-store/store/internal/adapters/rest/router"
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
// Метка path для GET-запросов сохраняется без идентификатора. Если путь не существует, добавляется метка со значением
// "non-existent path"
func (m *MiddlewareRequests) RequestsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var path string
		if r.Method == http.MethodGet {
			path = string([]rune(r.RequestURI)[0 : strings.LastIndex(r.RequestURI, "/")+1])
		} else {
			path = r.RequestURI
		}

		if router.IsExistPath(path) {
			m.metrics.HTTP.Requests.With(prometheus.Labels{metrics.PATH: path}).Inc()
		} else {
			m.metrics.HTTP.Requests.With(prometheus.Labels{metrics.PATH: "non-existent path"}).Inc()
		}
		next.ServeHTTP(rw, r)
	})
}
