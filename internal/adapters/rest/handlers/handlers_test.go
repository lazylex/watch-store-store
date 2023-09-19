package handlers

import (
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

func TestHandler_GetStockSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockService.NewMockInterface(ctrl)
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second))
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
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second))

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
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second))
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
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second))
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
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second))

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
	mux := router.New(&config.Config{Env: config.EnvironmentLocal}, New(service, nullLogger, 1*time.Second))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/api_v1/stock/amount/1", nil)

	service.EXPECT().GetAmountInStock(gomock.Any(), gomock.Any()).Times(1).Return(
		uint(0), repository.ErrNoRecord)

	mux.ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Fail()
	}
}
