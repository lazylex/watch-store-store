package handlers

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/lazylex/watch-store/store/internal/adapters/rest/router"
	"github.com/lazylex/watch-store/store/internal/config"
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

var nullLogger = logger.Null()

func nullSecure(next http.Handler) http.Handler {
	return next
}

func TestHandler_GetStockSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)
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
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/stock/9.9999", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_GetStockNoRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)
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
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)
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
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/stock/amount/1.0009", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_GetAmountInStockNoRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

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
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/price/9/1000", nil)

	service.EXPECT().ChangePriceInStock(gomock.Any(), dto.ArticleWithPriceDTO{Article: "9", Price: 1000})

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHandler_UpdatePriceInStockIncorrectArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/price/9.0090/1000", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_UpdatePriceInStockNegativePrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/price/9/-1000", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_UpdatePriceInStockIncorrectPrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/price/9/expensive-rich", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_UpdatePriceInStockTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/price/9/1000", nil)

	service.EXPECT().ChangePriceInStock(gomock.Any(),
		dto.ArticleWithPriceDTO{Article: "9", Price: 1000}).Times(1).Return(repository.ErrTimeout)

	mux.ServeHTTP(response, request)
	fmt.Println(response.Code)
	if response.Code != http.StatusRequestTimeout {
		t.Fail()
	}
}

func TestHandler_UpdateAmountInStockSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/amount/9/5", nil)

	service.EXPECT().ChangeAmountInStock(
		gomock.Any(), dto.ArticleWithAmountDTO{Article: "9", Amount: uint(5)}).Times(1).Return(nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHandler_UpdateAmountInStockIncorrectAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/amount/9/5five", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_UpdateAmountInStockIncorrectArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/amount/9.0999/5", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_UpdateAmountInStockTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/stock/amount/9/5", nil)

	service.EXPECT().ChangeAmountInStock(
		gomock.Any(), dto.ArticleWithAmountDTO{Article: "9", Amount: uint(5)}).Times(1).Return(repository.ErrTimeout)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusRequestTimeout {
		t.Fail()
	}
}

func TestHandler_AddToStockSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost,
		"/api/api_v1/stock/add/9/5/1000/test", nil)

	service.EXPECT().AddProductToStock(
		gomock.Any(), dto.NamedProductDTO{Article: "9", Amount: 5, Price: 1000, Name: "test"}).Times(1).Return(nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fail()
	}
}

func TestHandler_AddToStockIncorrectArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost,
		"/api/api_v1/stock/add/9.9999/5/1000/test", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_AddToStockIncorrectPrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost,
		"/api/api_v1/stock/add/9/5/too-many/test", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_AddToStockIncorrectAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost,
		"/api/api_v1/stock/add/9/five/1000/test", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_GetSoldAmountSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

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
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9.9999", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_GetSoldAmountTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

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
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9/2022-01-01/2023-09-28", nil)

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
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9/2024-01-01/2023-09-28", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_GetSoldAmountInTimePeriodTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9/2022-01-01/2023-09-28", nil)

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
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9/yesterday/2023-09-28", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_GetSoldAmountInTimePeriodIncorrectTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/sold/amount/9/2023-09-28/light-future", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_CancelReservationSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/reservation/cancel/9", nil)

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
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/reservation/cancel/nine", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_CancelReservationNegativeOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/reservation/cancel/-9", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_MakeReservationSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/api_v1/reservation/make?products[]=ca-f91w,2100,20&products[]=ca-aw-591,15000,36&order_number=9&status=0",
		nil)

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
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/api_v1/reservation/make?order_number=9&status=0",
		nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_MakeReservationNoStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/api_v1/reservation/make?products[]=ca-f91w,2100,20&products[]=ca-aw-591,15000,36&order_number=9",
		nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_MakeReservationNoOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/api_v1/reservation/make?products[]=ca-f91w,2100,20&products[]=ca-aw-591,15000,36&status=0",
		nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_MakeReservationIncorrectOrderData(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/api_v1/reservation/make?products[]=ca-f91w,2100,20&products[]=ca-aw-591,15000,36&status=4&order_number=5",
		nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_MakeLocalSaleSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost, "/api/api_v1/sale/make?products[]=ca-f91w,2100,20&products[]=ca-aw-591,15000,36",
		nil)

	service.EXPECT().MakeSale(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fail()
	}
}

func TestHandler_MakeLocalSaleNoProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/api_v1/sale/make", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_MakeLocalSaleTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost, "/api/api_v1/sale/make?products[]=ca-f91w,2100,20&products[]=ca-aw-591,15000,36",
		nil)

	service.EXPECT().MakeSale(gomock.Any(), gomock.Any()).Times(1).Return(repository.ErrTimeout)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusRequestTimeout {
		t.Fail()
	}
}

func TestHandler_MakeLocalSaleErrorData(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost, "/api/api_v1/sale/make?products[]=ca-f91w.2100.20&products[]=ca-aw-591.15000.36",
		nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_FinishOrderSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/reservation/finish/9", nil)

	service.EXPECT().FinishOrder(gomock.Any(), dto.OrderNumberDTO{OrderNumber: 9}).Times(1).Return(nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHandler_FinishOrderNegativeOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/reservation/finish/-9", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_FinishOrderNotIntOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/reservation/finish/nine", nil)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestHandler_FinishOrderTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second), nullSecure)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/api_v1/reservation/finish/9", nil)

	service.EXPECT().FinishOrder(gomock.Any(), dto.OrderNumberDTO{OrderNumber: 9}).Times(1).Return(repository.ErrTimeout)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusRequestTimeout {
		t.Fail()
	}
}
