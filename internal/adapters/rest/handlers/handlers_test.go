package handlers

import (
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/lazylex/watch-store/store/internal/adapters/rest/router"
	"github.com/lazylex/watch-store/store/internal/dto"
	"github.com/lazylex/watch-store/store/internal/logger"
	"github.com/lazylex/watch-store/store/internal/ports/repository"
	mockService "github.com/lazylex/watch-store/store/internal/ports/service/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandler_GetStockSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))
	service.EXPECT().GetStock(gomock.Any(), gomock.Any()).Times(1).Return(
		dto.NamedProductDTO{Name: "CASIO G-SHOCK DW-5600E-1V", Article: "1",
			Price: 7950, Amount: 22,
		}, nil)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/stock/1", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHandler_GetStockBadArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/stock/9.9999", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_GetStockNoRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))
	service.EXPECT().GetStock(gomock.Any(), gomock.Any()).Times(1).Return(
		dto.NamedProductDTO{}, repository.ErrNoRecord)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/stock/10000000000", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Fail()
	}
}

func TestHandler_GetAmountInStockSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))
	service.EXPECT().GetAmountInStock(gomock.Any(), gomock.Any()).Times(1).Return(
		uint(15), nil)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/stock/amount/1", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK || strings.Compare(response.Body.String(), "{\"amount\":15}\n") != 0 {
		t.Fail()
	}
}

func TestHandler_GetAmountInStockIncorrectArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/stock/amount/1.0009", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_GetAmountInStockNoRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/stock/amount/1", nil)

	service.EXPECT().GetAmountInStock(gomock.Any(), gomock.Any()).Times(1).Return(
		uint(0), repository.ErrNoRecord)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Fail()
	}
}

func TestHandler_UpdatePriceInStockSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.AddHandlers(chi.NewRouter(), New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/price", strings.NewReader("{\"Article\": \"9\", \"Price\": 1000}"))

	service.EXPECT().ChangePriceInStock(gomock.Any(), dto.ArticleWithPriceDTO{Article: "9", Price: 1000})

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHandler_UpdatePriceInStockIncorrectArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/price", strings.NewReader("{\"Article\": \"9.0090\", \"Price\": 1000}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_UpdatePriceInStockNegativePrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/price", strings.NewReader("{\"Article\": \"9\", \"Price\": -1000}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_UpdatePriceInStockIncorrectPrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/price", strings.NewReader("{\"Article\": \"9\", \"Price\": \"Expensive\"}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_UpdatePriceInStockTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/price", strings.NewReader("{\"Article\": \"9\", \"Price\": 1000}"))

	service.EXPECT().ChangePriceInStock(gomock.Any(),
		dto.ArticleWithPriceDTO{Article: "9", Price: 1000}).Times(1).Return(repository.ErrTimeout)

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusRequestTimeout {
		t.Fail()
	}
}

func TestHandler_UpdateAmountInStockSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/amount", strings.NewReader("{\"Article\": \"9\", \"Amount\": 5}"))

	service.EXPECT().ChangeAmountInStock(
		gomock.Any(), dto.ArticleWithAmountDTO{Article: "9", Amount: uint(5)}).Times(1).Return(nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHandler_UpdateAmountInStockIncorrectAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/amount", strings.NewReader("{\"Article\": \"9\", \"Amount\": \"five\"}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_UpdateAmountInStockIncorrectArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/amount", strings.NewReader("{\"Article\": \"9.0999\", \"Amount\": 5}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_UpdateAmountInStockTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/amount", strings.NewReader("{\"Article\": \"9\", \"Amount\": 5}"))

	service.EXPECT().ChangeAmountInStock(
		gomock.Any(), dto.ArticleWithAmountDTO{Article: "9", Amount: uint(5)}).Times(1).Return(repository.ErrTimeout)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusRequestTimeout {
		t.Fail()
	}
}

func TestHandler_AddToStockSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/api_v1/stock/add",
		strings.NewReader("{\"Article\": \"9\", \"Amount\": 5, \"Price\":1000, \"Name\":\"test\"}"))

	service.EXPECT().AddProductToStock(
		gomock.Any(), dto.NamedProductDTO{Article: "9", Amount: 5, Price: 1000, Name: "test"}).Times(1).Return(nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fail()
	}
}

func TestHandler_AddToStockIncorrectArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost,
		"/api/api_v1/stock/add",
		strings.NewReader("{\"Article\": \"9.9999\", \"Amount\": 5, \"Price\":1000, \"Name\":\"test\"}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_AddToStockIncorrectPrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost,
		"/api/api_v1/stock/add",
		strings.NewReader("{\"Article\": \"9\", \"Amount\": 5, \"Price\":\"too-many\", \"Name\":\"test\"}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_AddToStockIncorrectAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost,
		"/api/api_v1/stock/add",
		strings.NewReader("{\"Article\": \"9\", \"Amount\": \"five\", \"Price\":1000, \"Name\":\"test\"}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_GetSoldAmountSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9", nil)

	service.EXPECT().TotalSold(gomock.Any(), dto.ArticleDTO{Article: "9"}).Times(1).Return(uint(13), nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHandler_GetSoldAmountIncorrectArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9.9999", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_GetSoldAmountTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9", nil)

	service.EXPECT().TotalSold(gomock.Any(), dto.ArticleDTO{Article: "9"}).Times(1).Return(uint(0), repository.ErrTimeout)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusRequestTimeout {
		t.Fail()
	}
}

func TestHandler_GetSoldAmountInTimePeriodSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9?from=2022-01-01&to=2023-09-28", nil)

	service.EXPECT().TotalSoldInPeriod(
		gomock.Any(),
		dto.ArticleWithPeriodDTO{Article: "9",
			From: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			To:   time.Date(2023, 9, 28, 0, 0, 0, 0, time.UTC)}).Times(
		1).Return(uint(13), nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHandler_GetSoldAmountInTimePeriodIncorrectDateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9?from=2024-01-01&to=2023-09-28", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_GetSoldAmountInTimePeriodTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9?from=2022-01-01&to=2023-09-28", nil)

	service.EXPECT().TotalSoldInPeriod(
		gomock.Any(),
		dto.ArticleWithPeriodDTO{Article: "9",
			From: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			To:   time.Date(2023, 9, 28, 0, 0, 0, 0, time.UTC)}).Times(
		1).Return(uint(13), repository.ErrTimeout)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusRequestTimeout {
		t.Fail()
	}
}

func TestHandler_GetSoldAmountInTimePeriodIncorrectFrom(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9?from=yesterday&to=2023-09-28", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_GetSoldAmountInTimePeriodIncorrectTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9?from=2023-09-28&to=light-future", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_GetSoldAmountWithoutFrom(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9?to=2022-01-01", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_GetSoldAmountWithoutTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9?from=2022-01-01", nil)

	service.EXPECT().TotalSoldInPeriod(gomock.Any(), gomock.Any()).Times(1).Return(uint(13), nil)
	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHandler_CancelReservationSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/reservation/cancel",
		strings.NewReader("{\"order_number\": 9}"))

	service.EXPECT().CancelReservation(
		gomock.Any(),
		dto.OrderNumberDTO{OrderNumber: 9},
	).Times(1).Return(nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHandler_CancelReservationIncorrectOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/reservation/cancel",
		strings.NewReader("{\"order_number\": \"nine\"}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_CancelReservationNegativeOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/reservation/cancel",
		strings.NewReader("{\"order_number\": -9}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_MakeReservationSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/api_v1/reservation/make",
		strings.NewReader("{\"order_number\":13,\"state\":2,\"products\":[{\"article\":\"9\",\"price\":1330,\"amount\":6},{\"article\":\"1\",\"price\":3530,\"amount\":5}]}"))

	service.EXPECT().MakeReservation(
		gomock.Any(),
		gomock.Any(),
	).Times(1).Return(nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fail()
	}
}

func TestHandler_MakeReservationNoProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/api_v1/reservation/make",
		strings.NewReader("{\"order_number\":13,\"state\":1}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_MakeReservationNoStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/api_v1/reservation/make",
		strings.NewReader("{\"order_number\":9,\"products\":[{\"article\":\"9\",\"price\":1330,\"amount\":6},{\"article\":\"1\",\"price\":3530,\"amount\":5}]}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_MakeReservationIncorrectJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/api_v1/reservation/make",
		strings.NewReader("{\"order_number:9,\"products\":[{\"article\":\"9\",\"price\":1330,\"amount\":6},{\"article\":\"1\",\"price\":3530,\"amount\":5}]}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_MakeReservationNoOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/api_v1/reservation/make?products[]=ca-f91w,2100,20&products[]=ca-aw-591,15000,36&status=0",
		strings.NewReader("{\"state\":1,\"products\":[{\"article\":\"9\",\"price\":1330,\"amount\":6},{\"article\":\"1\",\"price\":3530,\"amount\":5}]}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_MakeReservationIncorrectOrderData(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/api_v1/reservation/make",
		strings.NewReader("{\"order_number\":5,\"state\":5,\"products\":[{\"article\":\"9\",\"price\":1330,\"amount\":6},{\"article\":\"1\",\"price\":3530,\"amount\":5}]}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_MakeLocalSaleSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost, "/api/api_v1/sale/make",
		strings.NewReader("[{\"article\":\"9\",\"price\":1330,\"amount\":6},{\"article\":\"1\",\"price\":3530,\"amount\":5}]"))

	service.EXPECT().MakeSale(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fail()
	}
}

func TestHandler_MakeLocalSaleNoProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/api_v1/sale/make", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_MakeLocalSaleTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost, "/api/api_v1/sale/make",
		strings.NewReader("[{\"article\":\"9\",\"price\":1330,\"amount\":6},{\"article\":\"1\",\"price\":3530,\"amount\":5}]"))

	service.EXPECT().MakeSale(gomock.Any(), gomock.Any()).Times(1).Return(repository.ErrTimeout)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusRequestTimeout {
		t.Fail()
	}
}

func TestHandler_MakeLocalSaleErrorData(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost, "/api/api_v1/sale/make",
		strings.NewReader("[{\"article\":\"9,\"pri\":1330,\"amt\":6},{\"article\":\"1\",\"price\":3530,\"amount\":5}]"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_FinishOrderSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/reservation/finish",
		strings.NewReader("{\"order_number\": 9}"))

	service.EXPECT().FinishOrder(gomock.Any(), dto.OrderNumberDTO{OrderNumber: 9}).Times(1).Return(nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHandler_FinishOrderNegativeOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/reservation/finish",
		strings.NewReader("{\"order_number\": -9}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_FinishOrderNotIntOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/reservation/finish",
		strings.NewReader("{\"order_number\": \"nine\"}"))

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_FinishOrderTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	mux := chi.NewRouter()
	service := mockService.NewMockInterface(ctrl)
	router.AddHandlers(mux, New(service, logger.Null(), time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/reservation/finish",
		strings.NewReader("{\"order_number\": 9}"))

	service.EXPECT().FinishOrder(gomock.Any(), dto.OrderNumberDTO{OrderNumber: 9}).Times(1).Return(repository.ErrTimeout)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusRequestTimeout {
		t.Fail()
	}
}
