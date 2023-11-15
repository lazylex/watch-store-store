package metrics

import (
	"errors"
	"fmt"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/internal/logger"
	httpMetrics "github.com/lazylex/watch-store/store/internal/ports/metrics/http"
	"github.com/lazylex/watch-store/store/internal/ports/metrics/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
)

const NAMESPACE = "store"

type Metrics struct {
	HTTP    httpMetrics.MetricsInterface
	Service service.MetricsInterface
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
		err                                                               error
		requests, canceledOrders, placedInternetOrders, placedLocalOrders *prometheus.CounterVec
		requestDuration                                                   *prometheus.HistogramVec
	)

	requests, err = createHTTPRequestsTotalMetric()
	if err != nil {
		return nil, err
	}

	requestDuration, err = createHTTPRequestDurationSecondsBucketMetric()
	if err != nil {
		return nil, err
	}

	canceledOrders, err = createCanceledOrdersTotalMetric()
	if err != nil {
		return nil, err
	}

	placedInternetOrders, err = createPlacedInternetOrdersTotalMetric()
	if err != nil {
		return nil, err
	}

	placedLocalOrders, err = createPlacedLocalOrdersTotalMetric()
	if err != nil {
		return nil, err
	}

	return &Metrics{
		Service: &Service{
			canceledOrders:       canceledOrders,
			placedLocalOrders:    placedLocalOrders,
			placedInternetOrders: placedInternetOrders},
		HTTP: &HTTP{requests: requests, duration: requestDuration},
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
