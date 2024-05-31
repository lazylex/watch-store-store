package server

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	restHandlers "github.com/lazylex/watch-store/store/internal/adapters/rest/handlers"
	"github.com/lazylex/watch-store/store/internal/adapters/rest/middlewares/jwt"
	requestMetrics "github.com/lazylex/watch-store/store/internal/adapters/rest/middlewares/request_metrics"
	"github.com/lazylex/watch-store/store/internal/adapters/rest/router"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/internal/helpers/constants/prefixes"
	"github.com/lazylex/watch-store/store/internal/logger"
	"github.com/lazylex/watch-store/store/internal/metrics"
	"github.com/lazylex/watch-store/store/internal/service"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Server struct {
	srv             *http.Server
	shutdownTimeout time.Duration
}

// serverErr добавляет к тексту ошибки префикс, указывающий на её принадлежность к DTO.
func serverErr(text string) error {
	return errors.New(prefixes.RestServerPrefix + text)
}

var (
	ErrNilConfigPointer  = serverErr("nil config pointer")
	ErrNilServicePointer = serverErr("nil domain service pointer")
	ErrNilMetricsPointer = serverErr("nil metrics pointer")
)

// MustCreate возвращает REST-сервер, который является оберткой над http.Server с настроенными middlewares и ручками. В
// качестве параметров передается адрес, таймауты, доменный сервис, метрики, окружение и строка безопасности для
// JWT-токена, если передано не локальное окружение для запуска.
func MustCreate(cfg *config.HttpServer, queryTimeout time.Duration,
	domainService *service.Service,
	metrics *metrics.Metrics,
	environment,
	signature string) *Server {
	handlers := restHandlers.New(domainService, queryTimeout)
	mux := chi.NewRouter()
	rm := requestMetrics.New(metrics)
	mux.Use(middleware.Recoverer, middleware.RequestID, rm.BeforeHandle, rm.AfterHandle)

	log := slog.With(logger.OPLabel, "adapters.rest.server.MustCreate")
	if cfg == nil {
		log.Error(ErrNilConfigPointer.Error())
		os.Exit(1)
	}
	if metrics == nil {
		log.Error(ErrNilMetricsPointer.Error())
		os.Exit(1)
	}
	if domainService == nil {
		log.Error(ErrNilServicePointer.Error())
		os.Exit(1)
	}

	if environment == config.EnvironmentLocal {
		mux.Use(middleware.Logger)
	} else {
		mux.Use(jwt.New([]byte(signature)).CheckJWT)
	}

	return &Server{
		srv: &http.Server{
			Handler:      router.AddHandlers(mux, handlers),
			Addr:         cfg.Address,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		shutdownTimeout: cfg.ShutdownTimeout,
	}
}

// MustRun производит запуск сервера в отдельной go-рутине. В случае ошибки останавливает работу приложения.
func (s *Server) MustRun() {
	log := slog.With(logger.OPLabel, "adapters.rest.server.MustRun")
	go func() {
		log.Info("start http server on " + s.srv.Addr)
		err := s.srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server startup error")
			os.Exit(1)
		}
	}()
}

// Shutdown производит остановку сервера.
func (s *Server) Shutdown() {
	log := slog.With(logger.OPLabel, "adapters.rest.server.Shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Error("failed to gracefully shutdown http server")
	} else {
		log.Info("gracefully shut down http server")
	}
}
