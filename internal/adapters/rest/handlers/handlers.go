package handlers

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/lazylex/watch-store/store/internal/adapters/rest/request"
	"github.com/lazylex/watch-store/store/internal/adapters/rest/response"
	"github.com/lazylex/watch-store/store/internal/domain/aggregates/reservation"
	"github.com/lazylex/watch-store/store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store/store/internal/dto"
	"github.com/lazylex/watch-store/store/internal/logger"
	"github.com/lazylex/watch-store/store/internal/ports/service"
	"log/slog"
	"net/http"
	"time"
)

type Handler struct {
	logger       *slog.Logger
	service      service.Interface
	queryTimeout time.Duration
}

// New конструктор хендлеров. Возвращает созданный обработчик *Handler
func New(service service.Interface, logger *slog.Logger, queryTimeout time.Duration) *Handler {
	return &Handler{logger: logger, service: service, queryTimeout: queryTimeout}
}

// injectRequestIDToCtx возвращает контекст с внедренным идентификатором запроса, для дальнейшего использования логгером
func injectRequestIDToCtx(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, logger.RequestId, middleware.GetReqID(r.Context()))
}

// GetStockRecord получение всех полей записи с переданным в пути запроса артикулом и возврат в
// формате JSON. Пример возвращаемых данных:
//
//	{
//	   "article": "CA-F91W",
//	   "name": "CASIO F-91W-1YEG",
//	   "price": 3490,
//	   "amount": 60
//	}
func (h *Handler) GetStockRecord(w http.ResponseWriter, r *http.Request) {
	var err error
	var art article.Article
	var stock dto.NamedProductDTO
	log := logger.AddPlaceAndRequestId(h.logger, "rest.handlers.GetStockRecord", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if art, err = request.GetArticleUsingChi(w, r, log); err != nil {
		return
	}

	transferObject := dto.ArticleDTO{Article: art}
	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	stock, err = h.service.GetStock(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	log.Info(fmt.Sprintf("requested stock record with article %s", art))

	render.JSON(w, r, stock)
}

// GetAmountInStock возвращает в формате JSON доступное для продажи количество товара с переданным в пути запроса
// артикулом. Пример возвращаемого значения:
// {
// "amount": 13
// }
func (h *Handler) GetAmountInStock(w http.ResponseWriter, r *http.Request) {
	var err error
	var art article.Article
	var amount uint
	log := logger.AddPlaceAndRequestId(h.logger, "rest.handlers.GetAmountInStock", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if art, err = request.GetArticleUsingChi(w, r, log); err != nil {
		return
	}

	transferObject := dto.ArticleDTO{Article: art}
	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	amount, err = h.service.GetAmountInStock(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	log.Info(fmt.Sprintf("requested amount in stock record with article %s", art))

	render.JSON(w, r, map[string]uint{request.Amount: amount})
}

// UpdatePriceInStock обновляет цену для товара. В пути запроса передается новая цена и артикул товара
func (h *Handler) UpdatePriceInStock(w http.ResponseWriter, r *http.Request) {
	var err error
	var art article.Article
	var price float64
	log := logger.AddPlaceAndRequestId(h.logger, "rest.handlers.UpdatePriceInStock", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if art, err = request.GetArticleUsingChi(w, r, log); err != nil {
		return
	}
	if price, err = request.GetPriceUsingChi(w, r, log); err != nil {
		return
	}

	transferObject := dto.ArticleWithPriceDTO{Article: art, Price: price}
	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	err = h.service.ChangePriceInStock(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	log.Info(fmt.Sprintf("price updated to %.2f in stock record with article %v", price, art))
}

// UpdateAmountInStock обновляет количество единиц товара. В пути запроса передаются новое количество и артикул товара
func (h *Handler) UpdateAmountInStock(w http.ResponseWriter, r *http.Request) {
	var err error
	var amount uint
	var art article.Article
	log := logger.AddPlaceAndRequestId(h.logger, "rest.handlers.UpdateAmountInStock", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if art, err = request.GetArticleUsingChi(w, r, log); err != nil {
		return
	}
	if amount, err = request.GetAmountUsingChi(w, r, log); err != nil {
		return
	}

	transferObject := dto.ArticleWithAmountDTO{Article: art, Amount: amount}
	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	err = h.service.ChangeAmountInStock(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	log.Info(fmt.Sprintf("amount udpaded to %d in stock record with article %v", amount, art))
}

// AddToStock добавляет новую запись о доступном товаре. В пути запроса передатся артикул, количество, цена и
// название товара
func (h *Handler) AddToStock(w http.ResponseWriter, r *http.Request) {
	var err error
	var name string
	var amount uint
	var price float64
	var art article.Article
	log := logger.AddPlaceAndRequestId(h.logger, "rest.handlers.CreateStock", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if art, err = request.GetArticleUsingChi(w, r, log); err != nil {
		return
	}
	if amount, err = request.GetAmountUsingChi(w, r, log); err != nil {
		return
	}
	if price, err = request.GetPriceUsingChi(w, r, log); err != nil {
		return
	}
	if name, err = request.GetNameUsingChi(w, r, log); err != nil {
		return
	}

	transferObject := dto.NamedProductDTO{Name: name, Article: art, Amount: amount, Price: price}
	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	err = h.service.AddProductToStock(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err == nil {
		w.WriteHeader(http.StatusCreated)
		log.Info(fmt.Sprintf("add product to stock with article %v, amount %d, price %.2f, name %s",
			art, amount, price, name))
	}
}

// GetSoldAmount возвращает общее количество проданного товара. В пути запроса передается артикул.
// Пример возвращаемого значения:
// {
// "amount": 13
// }
func (h *Handler) GetSoldAmount(w http.ResponseWriter, r *http.Request) {
	var err error
	var amount uint
	var art article.Article
	log := logger.AddPlaceAndRequestId(h.logger, "rest.handlers.GetSoldAmount", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if art, err = request.GetArticleUsingChi(w, r, log); err != nil {
		return
	}

	transferObject := dto.ArticleDTO{Article: art}
	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	amount, err = h.service.TotalSold(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	log.Info(fmt.Sprintf("requested total amount of sold with article %s", art))

	render.JSON(w, r, map[string]uint{request.Amount: amount})
}

// GetSoldAmountInTimePeriod возвращает общее количество проданного товара за период времени. В пути запроса передается
// артикул. Пример возвращаемого значения:
// {
// "amount": 13
// }
func (h *Handler) GetSoldAmountInTimePeriod(w http.ResponseWriter, r *http.Request) {
	var err error
	var amount uint
	var from, to time.Time
	var art article.Article
	log := logger.AddPlaceAndRequestId(h.logger, "rest.handlers.GetSoldAmountInTimePeriod", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if art, err = request.GetArticleUsingChi(w, r, log); err != nil {
		return
	}
	if from, err = request.GetFromUsingChi(w, r, log); err != nil {
		return
	}
	if to, err = request.GetToUsingChi(w, r, log); err != nil {
		return
	}

	transferObject := dto.ArticleWithPeriodDTO{Article: art, From: from, To: to}
	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	amount, err = h.service.TotalSoldInPeriod(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	log.Info(fmt.Sprintf("requested total amount of sold with article %s in time period", art))

	render.JSON(w, r, map[string]uint{request.Amount: amount})
}

// MakeReservation резервирует группу товаров под переданным номером заказа. В *http.Request передается номер заказа и
// массив резервируемых продуктов вида products[]=ca-f91w,2100,20&products[]=ca-aw-591,15000,36, где сперва идет
// артикул, затем цена и количество резервируемого товара. В случае удачного резервирования возвращается
// http.StatusCreated и производится запись в лог
func (h *Handler) MakeReservation(w http.ResponseWriter, r *http.Request) {
	var err error
	var status uint
	var products []dto.ProductDTO
	var order reservation.OrderNumber
	log := logger.AddPlaceAndRequestId(h.logger, "rest.handlers.MakeReservation", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if products, err = request.GetProductDTOs(w, r, log); err != nil {
		return
	}
	if order, err = request.GetOrderFromURLQuery(w, r, log); err != nil {
		return
	}
	if status, err = request.GetStatusFromURLQuery(w, r, log); err != nil {
		return
	}

	transferObject := dto.ReservationDTO{Products: products, OrderNumber: order, Date: time.Now(), State: status}
	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	err = h.service.MakeReservation(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err == nil {
		w.WriteHeader(http.StatusCreated)
		log.Info(fmt.Sprintf("succesfully saved order %d", order))
	}
}

// CancelReservation отменяет заказ с переданным в пути запроса номером заказа. Продукты из заказа возвращаются в
// продажу по актуальной на данный момент цене. В случае успешной отмены возвращается http.StatusOK и производится
// запись в лог
func (h *Handler) CancelReservation(w http.ResponseWriter, r *http.Request) {
	var err error
	var order reservation.OrderNumber
	log := logger.AddPlaceAndRequestId(h.logger, "rest.handlers.CancelReservation", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if order, err = request.GetOrderUsingChi(w, r, log); err != nil {
		return
	}

	transferObject := dto.OrderNumberDTO{OrderNumber: order}
	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	err = h.service.CancelReservation(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err == nil {
		log.Info(fmt.Sprintf("cancel order %d", order))
	}
}

// MakeLocalSale товар из доступного для продажи переносится в исторю продаж. В *http.Request передается массив
// реализуемых покупателю продуктов вида products[]=ca-f91w,2100,20&products[]=ca-aw-591,15000,36, где сперва идет
// артикул, затем цена и количество резервируемого товара. В случае удачного выполнения операции возвращается
// http.StatusOK и производится запись в лог
func (h *Handler) MakeLocalSale(w http.ResponseWriter, r *http.Request) {
	var err error
	var products []dto.ProductDTO
	log := logger.AddPlaceAndRequestId(h.logger, "rest.handlers.MakeLocalSale", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if products, err = request.GetProductDTOs(w, r, log); err != nil {
		return
	}

	err = h.service.MakeSale(injectRequestIDToCtx(ctx, r), products)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err == nil {
		var logString string
		for _, p := range products {
			logString += fmt.Sprintf("sold article: %s, amount: %d, price %.2f. ", p.Article, p.Amount, p.Price)
		}
		log.Info(logString)
	}
}

// FinishOrder отмечает заказ выполненым (отданым локальному покупателю или отправленным интернет-покупателю)
// и заносит зарезервированные продукты в историю проданных товаров
func (h *Handler) FinishOrder(w http.ResponseWriter, r *http.Request) {
	var err error
	var order reservation.OrderNumber
	log := logger.AddPlaceAndRequestId(h.logger, "rest.handlers.FinishOrder", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if order, err = request.GetOrderUsingChi(w, r, log); err != nil {
		return
	}
	transferObject := dto.OrderNumberDTO{OrderNumber: order}
	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	err = h.service.FinishOrder(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err == nil {
		log.Info(fmt.Sprintf("finish order %d", order))
	}
}
