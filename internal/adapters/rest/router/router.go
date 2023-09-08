package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	restHandles "github.com/lazylex/watch-store/store/internal/adapters/rest/handlers"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/internal/service"
	"log/slog"
)

const apiV1 = "/api/api_v1/"

// New возвращает роутер *chi.Mux для REST запросов
func New(cfg *config.Config, service *service.Service, logger *slog.Logger) *chi.Mux {
	router := chi.NewRouter()
	// ручки хардкорно создаются в роутере, а не передаются, как зависимость, потому что всё равно их логика считывания
	// параметров запросов местами зависима от chi, а значит с другим роутером они не совместимы без модификации
	handlers := restHandles.New(service, logger, cfg.QueryTimeout)
	if cfg.Env == config.EnvironmentLocal {
		router.Use(middleware.Logger)
	}

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	// TODO изменить пути на более соответствующие неймингу в сервисе
	router.Get(apiV1+"sale/{article}", handlers.GetStockRecord)
	router.Get(apiV1+"sale/amount/{article}", handlers.GetAmountInStock)
	router.Put(apiV1+"sale/amount/{article}/{amount}", handlers.UpdateAmountInStock)
	router.Put(apiV1+"sale/price/{article}/{price}", handlers.UpdatePriceInStock)
	router.Post(apiV1+"sale/add/{article}/{amount}/{price}/{name}", handlers.AddToStock)

	router.Get(apiV1+"sold/amount/{article}", handlers.GetSoldAmount)
	router.Get(apiV1+"sold/amount/{article}/{from}/{to}", handlers.GetSoldAmountInSpan)

	router.Post(apiV1+"sale/make", handlers.MakeLocalSale)
	router.Post(apiV1+"reservation/make", handlers.MakeReservation)
	router.Put(apiV1+"reservation/cancel/{order_number}", handlers.CancelReservation)

	return router
}
