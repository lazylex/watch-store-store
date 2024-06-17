// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/lazylex/watch-store-store/internal/dto"
)

// ключ для указания в контексте необходимости выполнения передаваемой функции, а не имитации её вызова
type ExecuteKey struct{}

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

// ConvertToCommonErr mocks base method.
func (m *MockInterface) ConvertToCommonErr(arg0 error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConvertToCommonErr", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConvertToCommonErr indicates an expected call of ConvertToCommonErr.
func (mr *MockInterfaceMockRecorder) ConvertToCommonErr(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConvertToCommonErr", reflect.TypeOf((*MockInterface)(nil).ConvertToCommonErr), arg0)
}

// CreateReservation mocks base method.
func (m *MockInterface) CreateReservation(arg0 context.Context, arg1 *dto.ReservationDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReservation", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateReservation indicates an expected call of CreateReservation.
func (mr *MockInterfaceMockRecorder) CreateReservation(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReservation", reflect.TypeOf((*MockInterface)(nil).CreateReservation), arg0, arg1)
}

// CreateSoldRecord mocks base method.
func (m *MockInterface) CreateSoldRecord(arg0 context.Context, arg1 *dto.SoldDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSoldRecord", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSoldRecord indicates an expected call of CreateSoldRecord.
func (mr *MockInterfaceMockRecorder) CreateSoldRecord(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSoldRecord", reflect.TypeOf((*MockInterface)(nil).CreateSoldRecord), arg0, arg1)
}

// CreateStock mocks base method.
func (m *MockInterface) CreateStock(arg0 context.Context, arg1 *dto.NamedProductDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStock", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateStock indicates an expected call of CreateStock.
func (mr *MockInterfaceMockRecorder) CreateStock(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStock", reflect.TypeOf((*MockInterface)(nil).CreateStock), arg0, arg1)
}

// DeleteReservation mocks base method.
func (m *MockInterface) DeleteReservation(arg0 context.Context, arg1 *dto.OrderNumberDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteReservation", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteReservation indicates an expected call of DeleteReservation.
func (mr *MockInterfaceMockRecorder) DeleteReservation(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteReservation", reflect.TypeOf((*MockInterface)(nil).DeleteReservation), arg0, arg1)
}

// ReadReservation mocks base method.
func (m *MockInterface) ReadReservation(arg0 context.Context, arg1 *dto.OrderNumberDTO) (dto.ReservationDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadReservation", arg0, arg1)
	ret0, _ := ret[0].(dto.ReservationDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadReservation indicates an expected call of ReadReservation.
func (mr *MockInterfaceMockRecorder) ReadReservation(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadReservation", reflect.TypeOf((*MockInterface)(nil).ReadReservation), arg0, arg1)
}

// ReadSoldAmount mocks base method.
func (m *MockInterface) ReadSoldAmount(arg0 context.Context, arg1 *dto.ArticleDTO) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadSoldAmount", arg0, arg1)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadSoldAmount indicates an expected call of ReadSoldAmount.
func (mr *MockInterfaceMockRecorder) ReadSoldAmount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadSoldAmount", reflect.TypeOf((*MockInterface)(nil).ReadSoldAmount), arg0, arg1)
}

// ReadSoldAmountInPeriod mocks base method.
func (m *MockInterface) ReadSoldAmountInPeriod(arg0 context.Context, arg1 *dto.ArticleWithPeriodDTO) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadSoldAmountInPeriod", arg0, arg1)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadSoldAmountInPeriod indicates an expected call of ReadSoldAmountInPeriod.
func (mr *MockInterfaceMockRecorder) ReadSoldAmountInPeriod(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadSoldAmountInPeriod", reflect.TypeOf((*MockInterface)(nil).ReadSoldAmountInPeriod), arg0, arg1)
}

// ReadSoldRecords mocks base method.
func (m *MockInterface) ReadSoldRecords(arg0 context.Context, arg1 *dto.ArticleDTO) ([]dto.SoldDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadSoldRecords", arg0, arg1)
	ret0, _ := ret[0].([]dto.SoldDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadSoldRecords indicates an expected call of ReadSoldRecords.
func (mr *MockInterfaceMockRecorder) ReadSoldRecords(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadSoldRecords", reflect.TypeOf((*MockInterface)(nil).ReadSoldRecords), arg0, arg1)
}

// ReadSoldRecordsInPeriod mocks base method.
func (m *MockInterface) ReadSoldRecordsInPeriod(arg0 context.Context, arg1 *dto.ArticleWithPeriodDTO) ([]dto.SoldDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadSoldRecordsInPeriod", arg0, arg1)
	ret0, _ := ret[0].([]dto.SoldDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadSoldRecordsInPeriod indicates an expected call of ReadSoldRecordsInPeriod.
func (mr *MockInterfaceMockRecorder) ReadSoldRecordsInPeriod(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadSoldRecordsInPeriod", reflect.TypeOf((*MockInterface)(nil).ReadSoldRecordsInPeriod), arg0, arg1)
}

// ReadStock mocks base method.
func (m *MockInterface) ReadStock(arg0 context.Context, arg1 *dto.ArticleDTO) (dto.NamedProductDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadStock", arg0, arg1)
	ret0, _ := ret[0].(dto.NamedProductDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadStock indicates an expected call of ReadStock.
func (mr *MockInterfaceMockRecorder) ReadStock(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadStock", reflect.TypeOf((*MockInterface)(nil).ReadStock), arg0, arg1)
}

// ReadStockAmount mocks base method.
func (m *MockInterface) ReadStockAmount(arg0 context.Context, arg1 *dto.ArticleDTO) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadStockAmount", arg0, arg1)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadStockAmount indicates an expected call of ReadStockAmount.
func (mr *MockInterfaceMockRecorder) ReadStockAmount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadStockAmount", reflect.TypeOf((*MockInterface)(nil).ReadStockAmount), arg0, arg1)
}

// ReadStockPrice mocks base method.
func (m *MockInterface) ReadStockPrice(arg0 context.Context, arg1 *dto.ArticleDTO) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadStockPrice", arg0, arg1)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadStockPrice indicates an expected call of ReadStockPrice.
func (mr *MockInterfaceMockRecorder) ReadStockPrice(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadStockPrice", reflect.TypeOf((*MockInterface)(nil).ReadStockPrice), arg0, arg1)
}

// UpdateReservation mocks base method.
func (m *MockInterface) UpdateReservation(arg0 context.Context, arg1 *dto.ReservationDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateReservation", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateReservation indicates an expected call of UpdateReservation.
func (mr *MockInterfaceMockRecorder) UpdateReservation(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateReservation", reflect.TypeOf((*MockInterface)(nil).UpdateReservation), arg0, arg1)
}

// UpdateStock mocks base method.
func (m *MockInterface) UpdateStock(arg0 context.Context, arg1 *dto.NamedProductDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStock", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStock indicates an expected call of UpdateStock.
func (mr *MockInterfaceMockRecorder) UpdateStock(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStock", reflect.TypeOf((*MockInterface)(nil).UpdateStock), arg0, arg1)
}

// UpdateStockAmount mocks base method.
func (m *MockInterface) UpdateStockAmount(arg0 context.Context, arg1 *dto.ArticleWithAmountDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStockAmount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStockAmount indicates an expected call of UpdateStockAmount.
func (mr *MockInterfaceMockRecorder) UpdateStockAmount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStockAmount", reflect.TypeOf((*MockInterface)(nil).UpdateStockAmount), arg0, arg1)
}

// UpdateStockPrice mocks base method.
func (m *MockInterface) UpdateStockPrice(arg0 context.Context, arg1 *dto.ArticleWithPriceDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStockPrice", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStockPrice indicates an expected call of UpdateStockPrice.
func (mr *MockInterfaceMockRecorder) UpdateStockPrice(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStockPrice", reflect.TypeOf((*MockInterface)(nil).UpdateStockPrice), arg0, arg1)
}

// WithinTransaction mocks base method.
func (m *MockInterface) WithinTransaction(arg0 context.Context, arg1 func(context.Context) error) error {
	// пришлось внести изменения в сгенерированный код, так как нужно тестировать логику, которую передают в функции arg1
	if ok := arg0.Value(ExecuteKey{}); ok != nil {
		return arg1(arg0)
	}
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithinTransaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// WithinTransaction indicates an expected call of WithinTransaction.
func (mr *MockInterfaceMockRecorder) WithinTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithinTransaction", reflect.TypeOf((*MockInterface)(nil).WithinTransaction), arg0, arg1)
}

// MockSQLDBInterface is a mock of SQLDBInterface interface.
type MockSQLDBInterface struct {
	ctrl     *gomock.Controller
	recorder *MockSQLDBInterfaceMockRecorder
}

// MockSQLDBInterfaceMockRecorder is the mock recorder for MockSQLDBInterface.
type MockSQLDBInterfaceMockRecorder struct {
	mock *MockSQLDBInterface
}

// NewMockSQLDBInterface creates a new mock instance.
func NewMockSQLDBInterface(ctrl *gomock.Controller) *MockSQLDBInterface {
	mock := &MockSQLDBInterface{ctrl: ctrl}
	mock.recorder = &MockSQLDBInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSQLDBInterface) EXPECT() *MockSQLDBInterfaceMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockSQLDBInterface) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockSQLDBInterfaceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSQLDBInterface)(nil).Close))
}

// DB mocks base method.
func (m *MockSQLDBInterface) DB() *sql.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DB")
	ret0, _ := ret[0].(*sql.DB)
	return ret0
}

// DB indicates an expected call of DB.
func (mr *MockSQLDBInterfaceMockRecorder) DB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DB", reflect.TypeOf((*MockSQLDBInterface)(nil).DB))
}
