// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/lazylex/watch-store-store/internal/dto"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// AddProductToStock mocks base method.
func (m *MockInterface) AddProductToStock(ctx context.Context, data dto.ArticlePriceNameAmount) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProductToStock", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddProductToStock indicates an expected call of AddProductToStock.
func (mr *MockInterfaceMockRecorder) AddProductToStock(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProductToStock", reflect.TypeOf((*MockInterface)(nil).AddProductToStock), ctx, data)
}

// AmountInStock mocks base method.
func (m *MockInterface) AmountInStock(ctx context.Context, data dto.Article) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AmountInStock", ctx, data)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AmountInStock indicates an expected call of AmountInStock.
func (mr *MockInterfaceMockRecorder) AmountInStock(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AmountInStock", reflect.TypeOf((*MockInterface)(nil).AmountInStock), ctx, data)
}

// CancelReservation mocks base method.
func (m *MockInterface) CancelReservation(ctx context.Context, data dto.Number) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelReservation", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelReservation indicates an expected call of CancelReservation.
func (mr *MockInterfaceMockRecorder) CancelReservation(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelReservation", reflect.TypeOf((*MockInterface)(nil).CancelReservation), ctx, data)
}

// ChangeAmountInStock mocks base method.
func (m *MockInterface) ChangeAmountInStock(ctx context.Context, data dto.ArticleAmount) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeAmountInStock", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeAmountInStock indicates an expected call of ChangeAmountInStock.
func (mr *MockInterfaceMockRecorder) ChangeAmountInStock(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeAmountInStock", reflect.TypeOf((*MockInterface)(nil).ChangeAmountInStock), ctx, data)
}

// ChangePriceInStock mocks base method.
func (m *MockInterface) ChangePriceInStock(ctx context.Context, data dto.ArticleWithPriceDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePriceInStock", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangePriceInStock indicates an expected call of ChangePriceInStock.
func (mr *MockInterfaceMockRecorder) ChangePriceInStock(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePriceInStock", reflect.TypeOf((*MockInterface)(nil).ChangePriceInStock), ctx, data)
}

// FinishOrder mocks base method.
func (m *MockInterface) FinishOrder(ctx context.Context, data dto.Number) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinishOrder", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// FinishOrder indicates an expected call of FinishOrder.
func (mr *MockInterfaceMockRecorder) FinishOrder(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinishOrder", reflect.TypeOf((*MockInterface)(nil).FinishOrder), ctx, data)
}

// MakeReservation mocks base method.
func (m *MockInterface) MakeReservation(ctx context.Context, data dto.ReservationDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeReservation", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// MakeReservation indicates an expected call of MakeReservation.
func (mr *MockInterfaceMockRecorder) MakeReservation(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeReservation", reflect.TypeOf((*MockInterface)(nil).MakeReservation), ctx, data)
}

// MakeSale mocks base method.
func (m *MockInterface) MakeSale(ctx context.Context, data []dto.ArticlePriceAmount) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeSale", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// MakeSale indicates an expected call of MakeSale.
func (mr *MockInterfaceMockRecorder) MakeSale(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeSale", reflect.TypeOf((*MockInterface)(nil).MakeSale), ctx, data)
}

// Stock mocks base method.
func (m *MockInterface) Stock(ctx context.Context, data dto.Article) (dto.ArticlePriceNameAmount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stock", ctx, data)
	ret0, _ := ret[0].(dto.ArticlePriceNameAmount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stock indicates an expected call of Stock.
func (mr *MockInterfaceMockRecorder) Stock(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stock", reflect.TypeOf((*MockInterface)(nil).Stock), ctx, data)
}

// TotalSold mocks base method.
func (m *MockInterface) TotalSold(ctx context.Context, data dto.Article) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TotalSold", ctx, data)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TotalSold indicates an expected call of TotalSold.
func (mr *MockInterfaceMockRecorder) TotalSold(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TotalSold", reflect.TypeOf((*MockInterface)(nil).TotalSold), ctx, data)
}

// TotalSoldInPeriod mocks base method.
func (m *MockInterface) TotalSoldInPeriod(ctx context.Context, data dto.ArticlePeriod) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TotalSoldInPeriod", ctx, data)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TotalSoldInPeriod indicates an expected call of TotalSoldInPeriod.
func (mr *MockInterfaceMockRecorder) TotalSoldInPeriod(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TotalSoldInPeriod", reflect.TypeOf((*MockInterface)(nil).TotalSoldInPeriod), ctx, data)
}
