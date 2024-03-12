package logger

import (
	"context"
	"github.com/go-chi/chi/middleware"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/pkg/colorlog"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type ContextKey int

const (
	RequestId     ContextKey = 0
	TxId          ContextKey = 1
	OPLabel                  = "op"
	RequestLabel             = "request_id"
	instanceLabel            = "instance"
	TxLabel                  = "tx_number"
)

// MustCreate возвращает экземпляр *slog.Logger или останавливает программу, если окружение environment указано неверно
func MustCreate(environment, instance string) *slog.Logger {
	var logger *slog.Logger
	switch environment {
	case config.EnvironmentLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvironmentDebug:
		logger = slog.New(colorlog.NewHandler(os.Stdout, &colorlog.Options{Level: slog.LevelDebug,
			TimeFormat: time.TimeOnly})).With(slog.String(instanceLabel, instance))
	case config.EnvironmentProduction:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		).With(slog.String(instanceLabel, instance))
	default:
		log.Fatal("program environment not set or it incorrect")
	}

	return logger
}

// AddPlaceAndRequestId добавляет в логгер место запуска и идентификатор запроса
func AddPlaceAndRequestId(log *slog.Logger, place string, r *http.Request) *slog.Logger {
	return log.With(
		slog.String(OPLabel, place), // operation place
		slog.String(RequestLabel, middleware.GetReqID(r.Context())),
	)
}

// LogWithCtxData извлекает, при наличии, из контекста идентификатор запроса и номер транзакции и добавляет в логгер
func LogWithCtxData(ctx context.Context, log *slog.Logger) *slog.Logger {
	if ctx.Value(RequestId) != nil {
		log = log.With(RequestLabel, ctx.Value(RequestId))
	}
	if ctx.Value(TxId) != nil {
		log = log.With(TxLabel, ctx.Value(TxId))
	}

	return log
}

// Null возвращает логгер, который пишет в никуда. Полезно для тестирования структур и функций, зависящих от логгера
func Null() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.NewFile(0, os.DevNull), nil))
}
