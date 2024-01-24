package metrics

import (
	"errors"
	"fmt"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/internal/helpers/constants/prefixes"
	internalLogger "github.com/lazylex/watch-store/store/internal/logger"
	httpMetrics "github.com/lazylex/watch-store/store/internal/ports/metrics/http"
	"github.com/lazylex/watch-store/store/internal/ports/metrics/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	standartLog "log"
	"log/slog"
	"net/http"
	"os"
)

const NAMESPACE = "store"

type Metrics struct {
	HTTP    httpMetrics.MetricsInterface
	Service service.MetricsInterface
}

// dtoErr добавляет к тексту ошибки префикс, указывающий на её принадлежность к DTO
func metricsErr(text string) error {
	return errors.New(prefixes.MetricsPrefix + text)
}

var (
	ErrNilConfigPointer = metricsErr("nil config pointer")
	ErrNilLoggerPointer = metricsErr("nil logger pointer")
)

// MustCreate возвращает метрики *Metrics или останавливает программу, если не удалось запустить http сервер для
// работы с Prometheus или занести метрики в регистр
func MustCreate(cfg *config.Prometheus, logger *slog.Logger) *Metrics {
	if logger == nil {
		standartLog.Fatal(ErrNilLoggerPointer.Error())
	}

	log := logger.With(slog.String(internalLogger.OPLabel, "metrics.MustCreate"))

	if cfg == nil {
		log.Error(ErrNilConfigPointer.Error())
		os.Exit(1)
	}

	var port = "9323"
	var url = "/metrics"

	if len(cfg.PrometheusPort) > 0 {
		port = cfg.PrometheusPort
	}

	if len(cfg.PrometheusMetricsURL) > 0 {
		url = cfg.PrometheusMetricsURL
	}

	startHTTP(url, port, logger)

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

	if log == nil {
		standartLog.Fatal(ErrNilLoggerPointer.Error())
	}

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
