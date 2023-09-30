package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lazylex/watch-store/store/internal/adapters/rest/middlewares/jwt"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/internal/ports/rest/handlers"
	"log/slog"
)

const apiV1 = "/api/api_v1/"

// New возвращает роутер *chi.Mux для REST запросов
func New(cfg *config.Config, logger *slog.Logger, handlers handlers.Interface) *chi.Mux {
	router := chi.NewRouter()

	if cfg.Env == config.EnvironmentLocal {
		router.Use(middleware.Logger)
	} else {
		// TODO: убрать хардхордно зашитую подпись для JWT токенов
		router.Use(jwt.New(logger, []byte("mySecretCombination")).CheckJWT)
	}

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)

	router.Get(apiV1+"stock/{article}", handlers.GetStockRecord)
	router.Get(apiV1+"stock/amount/{article}", handlers.GetAmountInStock)
	router.Put(apiV1+"stock/amount/{article}/{amount}", handlers.UpdateAmountInStock)
	router.Put(apiV1+"stock/price/{article}/{price}", handlers.UpdatePriceInStock)
	router.Post(apiV1+"stock/add/{article}/{amount}/{price}/{name}", handlers.AddToStock)

	router.Get(apiV1+"sold/amount/{article}", handlers.GetSoldAmount)
	router.Get(apiV1+"sold/amount/{article}/{from}/{to}", handlers.GetSoldAmountInTimePeriod)

	router.Post(apiV1+"sale/make", handlers.MakeLocalSale)

	router.Post(apiV1+"reservation/make", handlers.MakeReservation)
	router.Put(apiV1+"reservation/cancel/{order_number}", handlers.CancelReservation)
	router.Put(apiV1+"reservation/finish/{order_number}", handlers.FinishOrder)

	return router
}
