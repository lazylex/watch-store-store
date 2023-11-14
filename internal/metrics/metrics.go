package metrics

import (
	"errors"
	"fmt"
	"github.com/lazylex/watch-store/store/internal/adapters/rest/router"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/internal/helpers/constants/various"
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
	Duration *p.HistogramVec
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
	var (
		err             error
		requests        *p.CounterVec
		requestDuration *p.HistogramVec
	)

	requests, err = createHTTPRequestsTotalMetric()
	if err != nil {
		return nil, err
	}
	requestDuration, err = createHTTPRequestDurationSecondsBucketMetric()
	if err != nil {
		return nil, err
	}

	return &Metrics{
		Service: struct{}{},
		HTTP:    HTTP{Requests: requests, Duration: requestDuration},
	}, nil
}

// createHTTPRequestDurationSecondsBucketMetric создает и регистрирует метрику http_request_duration_seconds_bucket
func createHTTPRequestDurationSecondsBucketMetric() (*p.HistogramVec, error) {
	var err error
	requestDuration := p.NewHistogramVec(p.HistogramOpts{
		Namespace: NAMESPACE,
		Name:      "http_request_duration_seconds_bucket",
		Help:      "Duration of the request",
	}, []string{})
	if err = p.Register(requestDuration); err != nil {
		return nil, err
	}

	requestDuration.With(p.Labels{})

	return requestDuration, nil
}

// createHTTPRequestsTotalMetric создает и регистрирует метрику http_requests_total, являющуюся счетчиком http-запросов
func createHTTPRequestsTotalMetric() (*p.CounterVec, error) {
	var err error
	requests := p.NewCounterVec(p.CounterOpts{
		Name:      "http_requests_total",
		Namespace: NAMESPACE,
		Help:      "Count of http requests",
	}, []string{PATH})
	if err = p.Register(requests); err != nil {
		return nil, err
	}

	for _, path := range router.ExistentPaths() {
		requests.With(p.Labels{PATH: path})
	}
	requests.With(p.Labels{PATH: various.NonExistentPath})

	return requests, nil
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
