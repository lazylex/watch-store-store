package router

import (
	"github.com/go-chi/chi"
	"github.com/lazylex/watch-store/store/internal/ports/rest/handlers"
)

var paths map[string]string

// init регистрирует адреса в структуре paths. Ключами и значениями служат шаблоны пути, кроме адресов для GET-запросов.
// Для них ключом служит шаблон без передаваемого в фигурных скобках параметра. Регистрация адресов необходима для того,
// чтобы в метриках сохранялись только используемые в приложении пути.
func init() {
	paths = make(map[string]string, 10)
	paths["/api/api_v1/stock/"] = "/api/api_v1/stock/"
	paths["/api/api_v1/stock/amount/"] = "/api/api_v1/stock/amount/"
	paths["/api/api_v1/stock/amount"] = "/api/api_v1/stock/amount"
	paths["/api/api_v1/stock/price"] = "/api/api_v1/stock/price"
	paths["/api/api_v1/stock/add"] = "/api/api_v1/stock/add"
	paths["/api/api_v1/sold/amount/"] = "/api/api_v1/sold/amount/"
	paths["/api/api_v1/sale/make"] = "/api/api_v1/sale/make"
	paths["/api/api_v1/reservation/make"] = "/api/api_v1/reservation/make"
	paths["/api/api_v1/reservation/cancel"] = "/api/api_v1/reservation/cancel"
	paths["/api/api_v1/reservation/finish"] = "/api/api_v1/reservation/finish"
}

// ExistentPaths возвращает слайс доступных в приложении путей.
func ExistentPaths() []string {
	var existentPaths []string
	for k := range paths {
		existentPaths = append(existentPaths, k)
	}

	return existentPaths
}

// IsExistPath возвращает true, если в приложении используется передаваемый путь. Иначе - false.
func IsExistPath(path string) bool {
	if _, ok := paths[path]; ok {
		return true
	}
	return false
}

// AddHandlers возвращает роутер *chi.Mux с добавленными обработчиками REST запросов.
func AddHandlers(router *chi.Mux, handlers handlers.Interface) *chi.Mux {
	router.Get(paths["/api/api_v1/stock/"], handlers.StockRecord)
	router.Get(paths["/api/api_v1/stock/amount/"], handlers.AmountInStock)
	router.Put(paths["/api/api_v1/stock/amount"], handlers.UpdateAmountInStock)
	router.Put(paths["/api/api_v1/stock/price"], handlers.UpdatePriceInStock)
	router.Post(paths["/api/api_v1/stock/add"], handlers.AddToStock)

	router.Get(paths["/api/api_v1/sold/amount/"], handlers.SoldAmount)

	router.Post(paths["/api/api_v1/sale/make"], handlers.MakeLocalSale)

	router.Post(paths["/api/api_v1/reservation/make"], handlers.MakeReservation)
	router.Put(paths["/api/api_v1/reservation/cancel"], handlers.CancelReservation)
	router.Put(paths["/api/api_v1/reservation/finish"], handlers.FinishOrder)

	return router
}
