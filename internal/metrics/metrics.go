package metrics

import (
	"errors"
	"fmt"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/internal/logger"
	p "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
)

const NAMESPACE = "store"

type HTTP struct {
	Requests *p.CounterVec
}

type Service struct {
}

type Metrics struct {
	HTTP    HTTP
	Service Service
}

// MustCreate возвращает метрики *Metrics или останавливает программу, если не удалось запустить http сервер для
// работы с Prometheus или занести метрики в регистр
func MustCreate(cfg *config.Config, log *slog.Logger) *Metrics {
	var port = "9323"
	var url = "/metrics"

	log = log.With(slog.String(logger.OPLabel, "metrics.MustCreate"))

	if len(cfg.PrometheusPort) > 0 {
		port = cfg.PrometheusPort
	}

	if len(cfg.PrometheusMetricsURL) > 0 {
		url = cfg.PrometheusMetricsURL
	}

	startHTTP(url, port, log)

	metrics, err := registerMetrics()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	return metrics
}

// registerMetrics заносит метрики в регистр и возвращает их. При неудаче возвращает ошибку
func registerMetrics() (*Metrics, error) {
	var err error
	requests := p.NewCounterVec(p.CounterOpts{Name: "http_requests_total", Namespace: NAMESPACE}, []string{PATH})
	if err = p.Register(requests); err != nil {
		return nil, err
	}

	return &Metrics{
		Service: struct{}{},
		HTTP:    HTTP{Requests: requests},
	}, nil
}

// startHTTP запускает http сервер для связи с Prometheus на переданном в функцию порту и url. При неудаче выводит
// ошибку в лог и останавливает программу
func startHTTP(url, port string, log *slog.Logger) {
	go func() {
		mux := http.NewServeMux()

		mux.Handle(url, promhttp.Handler())
		log.Info(fmt.Sprintf(":%s%s ready for prometheus", port, url))
		err := http.ListenAndServe(":"+port, mux)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("can't start http server for prometheus")
			os.Exit(1)
		}
	}()
}
