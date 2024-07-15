package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/lazylex/watch-store-store/internal/adapters/rest/request"
	"github.com/lazylex/watch-store-store/internal/adapters/rest/response"
	"github.com/lazylex/watch-store-store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store-store/internal/dto"
	"github.com/lazylex/watch-store-store/internal/helpers/constants/various"
	"github.com/lazylex/watch-store-store/internal/logger"
	"github.com/lazylex/watch-store-store/internal/ports/service"
	"log/slog"
	"net/http"
	"time"
)

type Handler struct {
	service      service.Interface
	queryTimeout time.Duration
}

// New конструктор обработчиков. Возвращает созданный обработчик *Handler.
func New(service service.Interface, queryTimeout time.Duration) *Handler {
	return &Handler{service: service, queryTimeout: queryTimeout}
}

// injectRequestIDToCtx возвращает контекст с внедренным идентификатором запроса, для дальнейшей записи в лог.
func injectRequestIDToCtx(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, logger.RequestId, middleware.GetReqID(r.Context()))
}

// StockRecord получение всех полей записи с переданным в параметре запроса (article) артикулом и возврат в
// формате JSON. Пример возвращаемых данных:
//
//	{
//	   "article": "CA-F91W",
//	   "name": "CASIO F-91W-1YEG",
//	   "price": 3490,
//	   "amount": 60
//	}
func (h *Handler) StockRecord(w http.ResponseWriter, r *http.Request) {
	var err error
	var art article.Article
	var stock dto.ArticlePriceNameAmount
	log := logger.AddPlaceAndRequestId(slog.Default(), "rest.handlers.StockRecord", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	art = article.Article(r.FormValue(request.Article))

	transferObject := dto.Article{Article: art}
	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	stock, err = h.service.Stock(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	log.Info(fmt.Sprintf("requested stock record with article %s", art))

	render.JSON(w, r, stock)
}

// AmountInStock возвращает в формате JSON доступное для продажи количество товара с переданным параметром запроса
// (article) артикулом. Пример возвращаемого значения:
// {
// "amount": 13
// }
func (h *Handler) AmountInStock(w http.ResponseWriter, r *http.Request) {
	var err error
	var art article.Article
	var amount uint
	log := logger.AddPlaceAndRequestId(slog.Default(), "rest.handlers.AmountInStock", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	art = article.Article(r.FormValue(request.Article))

	transferObject := dto.Article{Article: art}
	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	amount, err = h.service.AmountInStock(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	log.Info(fmt.Sprintf("requested amount in stock record with article %s", art))

	render.JSON(w, r, map[string]uint{request.Amount: amount})
}

// UpdatePriceInStock обновляет цену для товара. В теле запроса передается новая цена и артикул товара в формате JSON
// Пример передаваемых данных:
//
// {"article":"3", "price":6759}
func (h *Handler) UpdatePriceInStock(w http.ResponseWriter, r *http.Request) {
	var err error
	var transferObject dto.ArticleWithPriceDTO

	log := logger.AddPlaceAndRequestId(slog.Default(), "rest.handlers.UpdatePriceInStock", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	err = json.NewDecoder(r.Body).Decode(&transferObject)
	if err != nil {
		response.WriteHeaderAndLogAboutBadRequest(w, log, err)
		return
	}

	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	err = h.service.ChangePriceInStock(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	log.Info(fmt.Sprintf("price updated to %.2f in stock record with article %v",
		transferObject.Price, transferObject.Article))
}

// UpdateAmountInStock обновляет количество единиц товара. В теле запроса передается новое количество и артикул товара в
// формате JSON. Пример передаваемых данных:
//
// {"article":"3", "amount":6759}
func (h *Handler) UpdateAmountInStock(w http.ResponseWriter, r *http.Request) {
	var err error
	var transferObject dto.ArticleAmount
	log := logger.AddPlaceAndRequestId(slog.Default(), "rest.handlers.UpdateAmountInStock", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	err = json.NewDecoder(r.Body).Decode(&transferObject)
	if err != nil {
		response.WriteHeaderAndLogAboutBadRequest(w, log, err)
		return
	}

	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	err = h.service.ChangeAmountInStock(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	log.Info(fmt.Sprintf("amount udpaded to %d in stock record with article %v",
		transferObject.Amount, transferObject.Article))
}

// AddToStock добавляет новую запись о доступном товаре. В теле запроса передается новое количество, артикул, цена и
// название товара в формате JSON. Пример передаваемых данных:
//
// {"article":"3", "amount":6759, "price": 16000, "name": "plum-bus"}
func (h *Handler) AddToStock(w http.ResponseWriter, r *http.Request) {
	var err error
	var transferObject dto.ArticlePriceNameAmount
	log := logger.AddPlaceAndRequestId(slog.Default(), "rest.handlers.CreateStock", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	err = json.NewDecoder(r.Body).Decode(&transferObject)
	if err != nil {
		response.WriteHeaderAndLogAboutBadRequest(w, log, err)
		return
	}

	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	err = h.service.AddProductToStock(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err == nil {
		w.WriteHeader(http.StatusCreated)
		log.Info(fmt.Sprintf("add product to stock with article %v, amount %d, price %.2f, name %s",
			transferObject.Article, transferObject.Amount, transferObject.Price, transferObject.Name))
	}
}

// SoldAmount возвращает общее количество проданного товара. В параметре запроса (article) передается артикул.
// Параметрами запроса опционально передаются даты from и to для указания временного диапазона. Если передать только
// параметр from, то в качестве параметра to будет текущая дата (определяется временем на сервере, где запущено
// приложение, а не БД). Если передан только параметр to, то возвращается ответ http.StatusBadRequest
// Пример возвращаемого значения:
// {
// "amount": 13
// }
func (h *Handler) SoldAmount(w http.ResponseWriter, r *http.Request) {
	var err error
	var amount uint
	var art article.Article
	log := logger.AddPlaceAndRequestId(slog.Default(), "rest.handlers.SoldAmount", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	art = article.Article(r.FormValue(request.Article))

	fromParam := r.FormValue(request.From)
	toParam := r.FormValue(request.To)

	if len(fromParam) == 0 && len(toParam) == 0 {
		transferObject := dto.Article{Article: art}
		err = transferObject.Validate()
		if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
			return
		}

		amount, err = h.service.TotalSold(injectRequestIDToCtx(ctx, r), transferObject)

	} else if len(fromParam) == 0 {
		response.WriteHeaderAndLogAboutBadRequest(w, log, request.ErrEmptyFromDate)
		return
	} else {
		var from, to time.Time

		if len(toParam) == 0 {
			to = time.Now()
		} else {
			if to, err = time.Parse(various.DateLayout, toParam); err != nil {
				response.WriteHeaderAndLogAboutBadRequest(w, log, request.ErrIncorrectDate)
				return
			}
		}

		if from, err = time.Parse(various.DateLayout, fromParam); err != nil {
			response.WriteHeaderAndLogAboutBadRequest(w, log, request.ErrIncorrectDate)
			return
		}

		transferObject := dto.ArticlePeriod{Article: art, From: from, To: to}
		err = transferObject.Validate()
		if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
			return
		}

		amount, err = h.service.TotalSoldInPeriod(injectRequestIDToCtx(ctx, r), transferObject)
	}

	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	log.Info(fmt.Sprintf("requested amount of sold with article %s", art))

	render.JSON(w, r, map[string]uint{request.Amount: amount})
}

// MakeReservation резервирует группу товаров под переданным номером заказа. В теле запроса передается номер заказа,
// статус резервирования (описание в internal/domain/aggregates/reservation/reservation.go) и массив резервируемых
// продуктов в формате JSON. В случае удачного резервирования возвращается http.StatusCreated и производится запись в
// лог. Пример передаваемых данных:
//
//	{
//		"order_number":13,
//		"state":1,
//		"products":[
//			{
//				"article" : "9",
//				"price" : 1330,
//				"amount":6
//			},
//			{
//				"article":"1",
//				"price":3530,
//				"amount":5
//			}
//		]
//	}
func (h *Handler) MakeReservation(w http.ResponseWriter, r *http.Request) {
	var err error
	var transferObject dto.ReservationDTO
	log := logger.AddPlaceAndRequestId(slog.Default(), "rest.handlers.MakeReservation", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	err = json.NewDecoder(r.Body).Decode(&transferObject)
	if err != nil {
		response.WriteHeaderAndLogAboutBadRequest(w, log, err)
		return
	}

	transferObject.Date = time.Now()
	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	err = h.service.MakeReservation(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err == nil {
		w.WriteHeader(http.StatusCreated)
		log.Info(fmt.Sprintf("succesfully saved order %d", transferObject.OrderNumber))
	}
}

// CancelReservation отменяет заказ с переданным в пути запроса номером заказа. Продукты из заказа возвращаются в
// продажу по актуальной на данный момент цене. В случае успешной отмены возвращается http.StatusOK и производится
// запись в лог. Данные в запросе передаются в теле в виде JSON. Например:
//
// {"order_number": 9}
func (h *Handler) CancelReservation(w http.ResponseWriter, r *http.Request) {
	var err error
	var transferObject dto.Number
	log := logger.AddPlaceAndRequestId(slog.Default(), "rest.handlers.CancelReservation", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	err = json.NewDecoder(r.Body).Decode(&transferObject)
	if err != nil {
		response.WriteHeaderAndLogAboutBadRequest(w, log, err)
		return
	}

	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	err = h.service.CancelReservation(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err == nil {
		log.Info(fmt.Sprintf("cancel order %d", transferObject.OrderNumber))
	}
}

// MakeLocalSale товар из доступного для продажи переносится в историю продаж. В случае удачного выполнения операции
// возвращается http.StatusOK и производится запись в лог. В теле запроса передается массив резервируемых продуктов
// в формате JSON. Пример передаваемых данных:
//
//	[
//		{
//			"article" : "9",
//			"price" : 1330,
//			"amount":6
//		},
//		{
//			"article":"1",
//			"price":3530,
//			"amount":5
//		}
//	]
func (h *Handler) MakeLocalSale(w http.ResponseWriter, r *http.Request) {
	var err error
	var products []dto.ArticlePriceAmount
	log := logger.AddPlaceAndRequestId(slog.Default(), "rest.handlers.MakeLocalSale", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	err = json.NewDecoder(r.Body).Decode(&products)
	if err != nil {
		response.WriteHeaderAndLogAboutBadRequest(w, log, err)
		return
	}

	err = h.service.MakeSale(injectRequestIDToCtx(ctx, r), products)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err == nil {
		var logString string
		for _, p := range products {
			logString += fmt.Sprintf("sold article: %s, amount: %d, price %.2f. ", p.Article, p.Amount, p.Price)
		}
		w.WriteHeader(http.StatusCreated)
		log.Info(logString)
	}
}

// FinishOrder отмечает заказ выполненным (отданным локальному покупателю или отправленным интернет-покупателю)
// и заносит зарезервированные продукты в историю проданных товаров. Данные в запросе передаются в теле в виде JSON.
// Например:
//
// {"order_number": 9}
func (h *Handler) FinishOrder(w http.ResponseWriter, r *http.Request) {
	var err error
	var transferObject dto.Number
	log := logger.AddPlaceAndRequestId(slog.Default(), "rest.handlers.FinishOrder", r)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	err = json.NewDecoder(r.Body).Decode(&transferObject)
	if err != nil {
		response.WriteHeaderAndLogAboutBadRequest(w, log, err)
		return
	}

	err = transferObject.Validate()
	if response.WriteHeaderAndLogAboutErr(w, log, err); err != nil {
		return
	}

	err = h.service.FinishOrder(injectRequestIDToCtx(ctx, r), transferObject)
	if response.WriteHeaderAndLogAboutErr(w, log, err); err == nil {
		log.Info(fmt.Sprintf("finish order %d", transferObject.OrderNumber))
	}
}
