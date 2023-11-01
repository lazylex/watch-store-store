package request

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/lazylex/watch-store/store/internal/adapters/rest/response"
	"github.com/lazylex/watch-store/store/internal/domain/aggregates/reservation"
	"github.com/lazylex/watch-store/store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store/store/internal/dto"
	"github.com/lazylex/watch-store/store/internal/helpers/constants/prefixes"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

const (
	Article     = "article"
	Products    = "products[]"
	Amount      = "amount"
	OrderNumber = "order_number"
	Status      = "status"
	From        = "from"
	To          = "to"
)

// requestErr добавляет к тексту ошибки префикс, указывающий на её принадлежность к запросу
func requestErr(text string) error {
	return errors.New(prefixes.RequestErrorsPrefix + text)
}

var ErrIncorrectOrder = requestErr("invalid order number passed")
var ErrIncorrectDate = requestErr("invalid date passed")
var ErrIncorrectReservationStatus = requestErr("invalid order status passed")
var ErrEmptyProductsData = requestErr("product data not sent")
var ErrEmptyFromDate = requestErr("no 'from' date in request")

// GetArticleUsingChi возвращает артукул продукта
func GetArticleUsingChi(r *http.Request) article.Article {
	return article.Article(chi.URLParam(r, Article))
}

// GetStatusFromURLQuery возвращает статус заказа, если он есть в запросе. При отсутствии статуса в заголовок ответа
// записывается http.StatusBadRequest, и возвращается ошибка ErrIncorrectReservationStatus, которая так же записывается
// в лог
func GetStatusFromURLQuery(w http.ResponseWriter, r *http.Request, logger *slog.Logger) (uint, error) {
	parseUint, err := strconv.ParseUint(r.URL.Query().Get(Status), 10, 64)
	if err != nil {
		response.WriteHeaderAndLogAboutBadRequest(w, logger, ErrIncorrectReservationStatus)
	}

	return uint(parseUint), err
}

// GetOrderUsingChi возвращает номер заказа, если он есть в запросе. При отсутствии номера в заголовок ответа
// записывается http.StatusBadRequest, и возвращается ошибка ErrIncorrectOrder, которая так же записывается в лог
func GetOrderUsingChi(w http.ResponseWriter, r *http.Request, logger *slog.Logger) (reservation.OrderNumber, error) {
	order, err := strconv.Atoi(chi.URLParam(r, OrderNumber))
	if err != nil {
		response.WriteHeaderAndLogAboutBadRequest(w, logger, ErrIncorrectOrder)
	}
	return reservation.OrderNumber(order), err
}

// GetOrderFromURLQuery возвращает номер заказа, если он есть в запросе. При отсутствии номера в заголовок ответа
// записывается http.StatusBadRequest, и возвращается ошибка ErrIncorrectOrder, которая так же записывается в лог
func GetOrderFromURLQuery(w http.ResponseWriter, r *http.Request, logger *slog.Logger) (reservation.OrderNumber, error) {
	order, err := strconv.Atoi(r.URL.Query().Get(OrderNumber))
	if err != nil {
		response.WriteHeaderAndLogAboutBadRequest(w, logger, ErrIncorrectOrder)
	}
	return reservation.OrderNumber(order), err
}

// GetProductDTOs возращает []dto.ProductDTO, полученный из аргументов запроса, содержащего данные о продукте в таком
// виде: products[]=ca-f91w,2100,20&products[]=ca-aw-591,15000,36. Если данных не обнаружено, в заголовок ответа
// записывается http.StatusBadRequest, возвращается пустой слайс и ошибка ErrEmptyProductsData, которая так же
// записывается в лог
func GetProductDTOs(w http.ResponseWriter, r *http.Request, logger *slog.Logger) ([]dto.ProductDTO, error) {
	var products []dto.ProductDTO
	articlePosition, pricePosition, amountPosition := 0, 1, 2
	for _, productsData := range r.URL.Query()[Products] {
		productData := strings.Split(productsData, ",")
		if len(productData) != 3 {
			continue
		}

		price, err := strconv.ParseFloat(productData[pricePosition], 64)
		if err != nil {
			continue
		}

		amount, err := strconv.ParseUint(productData[amountPosition], 10, 64)
		if err != nil {
			continue
		}

		art := article.Article(productData[articlePosition])

		productDTO := dto.ProductDTO{Article: art, Price: price, Amount: uint(amount)}
		if productDTO.Validate() == nil {
			products = append(products, productDTO)
		}
	}

	if len(products) == 0 {
		response.WriteHeaderAndLogAboutErr(w, logger, ErrEmptyProductsData)
		return products, ErrEmptyProductsData
	}

	return products, nil
}
