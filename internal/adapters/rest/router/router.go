package router

import (
	"github.com/go-chi/chi"
	"github.com/lazylex/watch-store/store/internal/ports/rest/handlers"
)

const apiV1 = "/api/api_v1/"

// AddHandlers возвращает роутер *chi.Mux с добавленными хендлерами REST запросов
func AddHandlers(router *chi.Mux, handlers handlers.Interface) *chi.Mux {
	router.Get(apiV1+"stock/{article}", handlers.GetStockRecord)
	router.Get(apiV1+"stock/amount/{article}", handlers.GetAmountInStock)
	router.Put(apiV1+"stock/amount", handlers.UpdateAmountInStock)
	router.Put(apiV1+"stock/price", handlers.UpdatePriceInStock)
	router.Post(apiV1+"stock/add/{article}/{amount}/{price}/{name}", handlers.AddToStock)

	router.Get(apiV1+"sold/amount/{article}", handlers.GetSoldAmount)
	router.Get(apiV1+"sold/amount/{article}/{from}/{to}", handlers.GetSoldAmountInTimePeriod)

	router.Post(apiV1+"sale/make", handlers.MakeLocalSale)

	router.Post(apiV1+"reservation/make", handlers.MakeReservation)
	router.Put(apiV1+"reservation/cancel/{order_number}", handlers.CancelReservation)
	router.Put(apiV1+"reservation/finish/{order_number}", handlers.FinishOrder)

	return router
}
