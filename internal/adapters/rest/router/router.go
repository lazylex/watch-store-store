package router

import (
	"github.com/go-chi/chi"
	"github.com/lazylex/watch-store/store/internal/adapters/rest/handlers"
	"net/http"
)

var paths []string

const (
	apiApiV1Stock             = "/api/api_v1/stock/"
	apiApiV1StockAmountGet    = "/api/api_v1/stock/amount/"
	apiApiV1StockAmountUpdate = "/api/api_v1/stock/amount"
	apiApiV1StockPrice        = "/api/api_v1/stock/price"
	apiApiV1StockAdd          = "/api/api_v1/stock/add"
	apiApiV1SoldAmount        = "/api/api_v1/sold/amount/"
	apiApiV1SaleMake          = "/api/api_v1/sale/make"
	apiApiV1ReservationMake   = "/api/api_v1/reservation/make"
	apiApiV1ReservationCancel = "/api/api_v1/reservation/cancel"
	apiApiV1ReservationFinish = "/api/api_v1/reservation/finish"
)

const (
	receiveProductData                 = "получать данные о товаре"
	getQuantityOfGoodsAvailableForSale = "получать доступное для продажи количество товара"
	updateProductQuantity              = "обновлять количество товара"
	updateProductPrice                 = "обновлять цену товара"
	addProductEntry                    = "добавлять запись о товаре"
	getTotalQuantityOfGoodsSold        = "получать общее количество проданного товара"
	carryOutLocalSales                 = "осуществлять локальную продажу"
	reserveGoods                       = "резервировать товар"
	cancelReservation                  = "отменять резервирование"
	completeSaleOrShipment             = "завершать продажу/отправку"
)

func init() {
	paths = []string{
		apiApiV1Stock,
		apiApiV1StockAmountGet,
		apiApiV1StockAmountUpdate,
		apiApiV1StockPrice,
		apiApiV1StockAdd,
		apiApiV1SoldAmount,
		apiApiV1SaleMake,
		apiApiV1ReservationMake,
		apiApiV1ReservationCancel,
		apiApiV1ReservationFinish,
	}
}

type Router struct {
	handlers *handlers.Handler
	routes   []Route
	mux      *chi.Mux
}

type Route struct {
	Path       string
	Method     string
	Permission string
	Handler    func(w http.ResponseWriter, r *http.Request)
}

// Mux возвращает мультиплексор.
func (r *Router) Mux() *chi.Mux {
	return r.mux
}

// MustCreate конструктор для маршрутизатора.
func MustCreate(handlers *handlers.Handler) *Router {
	r := Router{handlers: handlers}
	r.mux = chi.NewRouter()
	r.registerRoutes()
	return &r
}

// registerRoutes сохраняет в памяти пути, обработчики, методы доступа и необходимые для выполнения обработчика пути.
// Прикрепление путей к обработчикам в мультиплексоре осуществляется функцией AssignPathsToHandlers.
func (r *Router) registerRoutes() {
	r.routes = []Route{
		{
			Path:       apiApiV1Stock,
			Method:     http.MethodGet,
			Permission: receiveProductData,
			Handler:    r.handlers.StockRecord,
		},
		{
			Path:       apiApiV1StockAmountGet,
			Method:     http.MethodGet,
			Permission: getQuantityOfGoodsAvailableForSale,
			Handler:    r.handlers.AmountInStock,
		},
		{
			Path:       apiApiV1StockAmountUpdate,
			Method:     http.MethodPut,
			Permission: updateProductQuantity,
			Handler:    r.handlers.UpdateAmountInStock,
		},
		{
			Path:       apiApiV1StockPrice,
			Method:     http.MethodPut,
			Permission: updateProductPrice,
			Handler:    r.handlers.UpdatePriceInStock,
		},
		{
			Path:       apiApiV1StockAdd,
			Method:     http.MethodPost,
			Permission: addProductEntry,
			Handler:    r.handlers.AddToStock,
		},
		{
			Path:       apiApiV1SoldAmount,
			Method:     http.MethodGet,
			Permission: getTotalQuantityOfGoodsSold,
			Handler:    r.handlers.SoldAmount,
		},
		{
			Path:       apiApiV1SaleMake,
			Method:     http.MethodPost,
			Permission: carryOutLocalSales,
			Handler:    r.handlers.MakeLocalSale,
		},
		{
			Path:       apiApiV1ReservationMake,
			Method:     http.MethodPost,
			Permission: reserveGoods,
			Handler:    r.handlers.MakeReservation,
		},
		{
			Path:       apiApiV1ReservationCancel,
			Method:     http.MethodPut,
			Permission: cancelReservation,
			Handler:    r.handlers.CancelReservation,
		},
		{
			Path:       apiApiV1ReservationFinish,
			Method:     http.MethodPut,
			Permission: completeSaleOrShipment,
			Handler:    r.handlers.FinishOrder,
		},
	}
}

// AssignPathsToHandlers прикрепляет обработчики к зарегистрированным на них путям.
func (r *Router) AssignPathsToHandlers() {
	for _, route := range r.routes {
		switch route.Method {
		case http.MethodGet:
			r.mux.Get(route.Path, route.Handler)
		case http.MethodPost:
			r.mux.Post(route.Path, route.Handler)
		case http.MethodPut:
			r.mux.Put(route.Path, route.Handler)
		case http.MethodPatch:
			r.mux.Patch(route.Path, route.Handler)
		case http.MethodDelete:
			r.mux.Delete(route.Path, route.Handler)
		}
	}
}

// IsExistPath возвращает true, если в приложении существует передаваемый путь. Иначе - false.
func IsExistPath(path string) bool {
	for _, p := range ExistentPaths() {
		if p == path {
			return true
		}
	}

	return false
}

// ExistentPaths возвращает список существующих путей.
func ExistentPaths() []string {
	return paths
}
