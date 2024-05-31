package request_metrics

import (
	"github.com/lazylex/watch-store/store/internal/adapters/rest/router"
	"github.com/lazylex/watch-store/store/internal/helpers/constants/various"
	"github.com/lazylex/watch-store/store/internal/metrics"
	"net/http"
	"time"
)

type MiddlewareRequests struct {
	metrics *metrics.Metrics
}

// New конструктор прослойки для http-запросов.
func New(metrics *metrics.Metrics) *MiddlewareRequests {
	return &MiddlewareRequests{metrics: metrics}
}

// BeforeHandle - middleware для метрик, подсчитываемых до выполнения обработчика запроса.
func (m *MiddlewareRequests) BeforeHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		m.requestsInc(r)
		next.ServeHTTP(rw, r)
	})
}

// AfterHandle - middleware для метрик, подсчитываемых после выполнения обработчика запроса.
func (m *MiddlewareRequests) AfterHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			duration := float64(time.Now().UnixMilli()-start.UnixMilli()) * 0.001
			m.metrics.HTTP.RequestsDurationObserve(duration)
		}()

		next.ServeHTTP(rw, r)
	})
}

// requestsInc увеличивает счетчик http-запросов к приложению. Добавляется метка path, содержащая путь запроса.
// Если путь не существует, добавляется метка со значением "non-existent path".
func (m *MiddlewareRequests) requestsInc(r *http.Request) {
	if router.IsExistPath(r.URL.Path) {
		m.metrics.HTTP.RequestsTotalInc(map[string]string{metrics.PATH: r.URL.Path})
	} else {
		m.metrics.HTTP.RequestsTotalInc(map[string]string{metrics.PATH: various.NonExistentPath})
	}
}
