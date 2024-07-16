package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lazylex/watch-store-store/internal/domain/aggregates/reservation"
	"github.com/lazylex/watch-store-store/internal/dto"
	"github.com/lazylex/watch-store-store/internal/metrics"
	mockService "github.com/lazylex/watch-store-store/internal/ports/metrics/service/mocks"
	"github.com/lazylex/watch-store-store/internal/ports/repository"
	mockrepository "github.com/lazylex/watch-store-store/internal/ports/repository/mocks"
	"github.com/lazylex/watch-store-store/internal/ports/service"
	"os"
	"os/exec"
	"time"

	"testing"
)

// ключ для указания в контексте необходимости выполнения передаваемой функции, а не имитации её вызова
type ExecuteKey struct{}

// withMockRepo позволяет подключать мок репозитория в функции service.New, чтобы покрыть её тестами
func withMockRepo(mr repository.Interface) Option {
	return func(s *Service) {
		s.Repository = mr
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	if os.Getenv("BE_CRASHER") == "1" {
		New()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestNew")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	var e *exec.ExitError
	if errors.As(err, &e) && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestService_AddProductToStockCorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticlePriceNameAmount{Name: "test_correct", Article: "test-9", Price: 110, Amount: 10}
	s := New(withMockRepo(mockRepo), WithMetrics(nil))

	mockRepo.EXPECT().CreateStock(context.Background(), &data).Times(1).Return(nil)

	err := s.AddProductToStock(context.Background(), data)
	if err != nil {
		t.Fail()
	}
}

func TestService_AddProductToStockIncorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticlePriceNameAmount{Name: "test_incorrect", Article: "test-9.9999", Price: -110, Amount: 10}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().CreateStock(context.Background(), &data).Times(0)

	err := s.AddProductToStock(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_AddProductToStockDuplicateArticle(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticlePriceNameAmount{Name: "test_correct", Article: "test-9", Price: 110, Amount: 10}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().CreateStock(context.Background(), &data).Times(1).Return(nil)

	err := s.AddProductToStock(context.Background(), data)
	if err != nil {
		t.Fail()
	}

	mockRepo.EXPECT().CreateStock(context.Background(), &data).Times(1).Return(errors.New("already exist"))

	err = s.AddProductToStock(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_ChangePriceInStockCorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticlePrice{Article: "test-9", Price: 10}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().ReadStock(gomock.Any(), gomock.Any()).Times(1)
	mockRepo.EXPECT().UpdateStockPrice(context.Background(), &data).Times(1).Return(nil)

	err := s.ChangePriceInStock(context.Background(), data)
	if err != nil {
		t.Fail()
	}
}

func TestService_ChangePriceInStockIncorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticlePrice{Article: "test-9", Price: -10}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().UpdateStockPrice(context.Background(), &data).Times(0)

	err := s.ChangePriceInStock(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_ChangePriceInStockNoInRepo(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticlePrice{Article: "test-9", Price: 100}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().ReadStock(context.Background(), &dto.Article{Article: "test-9"}).Times(1).
		Return(dto.ArticlePriceNameAmount{}, errors.New("no in stock"))
	mockRepo.EXPECT().UpdateStockPrice(context.Background(), &data).Times(0)

	err := s.ChangePriceInStock(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_GetStockCorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.Article{Article: "test-9"}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().ReadStock(context.Background(), &data).Times(1).Return(dto.ArticlePriceNameAmount{
		Name: "test-9", Article: "test-9", Price: 110, Amount: 10}, nil)

	_, err := s.Stock(context.Background(), data)
	if err != nil {
		t.Fail()
	}
}

func TestService_GetStockIncorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.Article{Article: "test-9.9999"}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().ReadStock(context.Background(), &data).Times(0)

	_, err := s.Stock(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_GetStockNoRecord(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.Article{Article: "test-9"}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().ReadStock(context.Background(), &data).Times(1).Return(dto.ArticlePriceNameAmount{},
		repository.ErrNoRecord)

	_, err := s.Stock(context.Background(), data)
	if !errors.Is(err, repository.ErrNoRecord) {
		t.Fail()
	}
}

func TestService_ChangeAmountInStockCorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleAmount{Article: "test-9", Amount: 10}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().UpdateStockAmount(context.Background(), &data).Times(1).Return(nil)

	err := s.ChangeAmountInStock(context.Background(), data)
	if err != nil {
		t.Fail()
	}
}

func TestService_ChangeAmountInStockIncorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleAmount{Article: "test-9.9999", Amount: 10}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().UpdateStockAmount(context.Background(), &data).Times(0)

	err := s.ChangeAmountInStock(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_ChangeAmountInStockNoRecord(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleAmount{Article: "test-9", Amount: 10}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().UpdateStockAmount(context.Background(), &data).Times(1).Return(repository.ErrNoRecord)

	err := s.ChangeAmountInStock(context.Background(), data)
	if !errors.Is(err, repository.ErrNoRecord) {
		t.Fail()
	}
}

func TestService_GetAmountInStockCorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.Article{Article: "test-9"}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().ReadStockAmount(context.Background(), &data).Times(1).Return(uint(5), nil)

	_, err := s.AmountInStock(context.Background(), data)
	if err != nil {
		t.Fail()
	}
}

func TestService_GetAmountInStockIncorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.Article{Article: "test-9.9999"}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().ReadStockAmount(context.Background(), &data).Times(0)

	_, err := s.AmountInStock(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_GetAmountInStockNoRecord(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.Article{Article: "test-9"}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().ReadStockAmount(context.Background(), &data).Times(1).Return(uint(0), repository.ErrNoRecord)

	_, err := s.AmountInStock(context.Background(), data)
	if !errors.Is(err, repository.ErrNoRecord) {
		t.Fail()
	}
}

func TestService_TotalSoldCorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.Article{Article: "test-9"}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().ReadSoldAmount(context.Background(), &data).Times(1).Return(uint(5), nil)

	_, err := s.TotalSold(context.Background(), data)
	if err != nil {
		t.Fail()
	}
}

func TestService_TotalSoldIncorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.Article{Article: "test-9.9999"}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().ReadSoldAmount(context.Background(), &data).Times(0)

	_, err := s.TotalSold(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_TotalSoldTimeout(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.Article{Article: "test-9"}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().ReadSoldAmount(context.Background(), &data).Times(1).Return(uint(0), repository.ErrTimeout)

	_, err := s.TotalSold(context.Background(), data)
	if !errors.Is(err, repository.ErrTimeout) {
		t.Fail()
	}
}

func TestService_TotalSoldInPeriodCorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleFromTo{
		Article: "test_9",
		From:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		To:      time.Now(),
	}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().ReadSoldAmountInPeriod(context.Background(), &data).Times(1).Return(uint(5), nil)

	_, err := s.TotalSoldInPeriod(context.Background(), data)
	if err != nil {
		t.Fail()
	}
}

func TestService_TotalSoldInPeriodIncorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleFromTo{Article: "test_9.9999", From: time.Now(), To: time.Now()}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().ReadSoldAmountInPeriod(context.Background(), &data).Times(0)

	_, err := s.TotalSoldInPeriod(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_TotalSoldInPeriodTimeout(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleFromTo{
		Article: "test_9",
		From:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		To:      time.Now(),
	}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().ReadSoldAmountInPeriod(context.Background(), &data).Times(1).Return(
		uint(0), repository.ErrTimeout)

	_, err := s.TotalSoldInPeriod(context.Background(), data)
	if !errors.Is(err, repository.ErrTimeout) {
		t.Fail()
	}
}

func TestService_MakeReservationIncorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.NumberDateStateProducts{
		Products:    []dto.ArticlePriceAmount{{Article: "test-9", Amount: 1, Price: 698}},
		OrderNumber: reservation.MaxCashRegisterNumber + 1,
		Date:        time.Now(),
		State:       reservation.NewForCashRegister,
	}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().WithinTransaction(context.Background(), gomock.Any()).Times(0)

	err := s.MakeReservation(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_MakeReservationSuccess(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.NumberDateStateProducts{
		Products:    []dto.ArticlePriceAmount{{Article: "test-9", Amount: 1, Price: 698}},
		OrderNumber: reservation.MaxCashRegisterNumber,
		Date:        time.Now(),
		State:       reservation.NewForCashRegister,
	}
	s := Service{Repository: mockRepo}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(uint(5), nil)
	mockRepo.EXPECT().UpdateStockAmount(ctx,
		&dto.ArticleAmount{Article: "test-9", Amount: uint(4)}).Times(1).Return(nil)
	mockRepo.EXPECT().CreateReservation(ctx, &data).Times(1).Return(nil)

	err := s.MakeReservation(ctx, data)
	if err != nil {
		t.Fail()
	}
}

func TestService_MakeReservationForInternetSuccess(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	mockServiceMetrics := mockService.NewMockMetricsInterface(ctrl)
	data := dto.NumberDateStateProducts{
		Products:    []dto.ArticlePriceAmount{{Article: "test-9", Amount: 1, Price: 698}},
		OrderNumber: reservation.MaxCashRegisterNumber + 1,
		Date:        time.Now(),
		State:       reservation.NewForInternetCustomer,
	}
	s := Service{Repository: mockRepo,
		Metrics: &metrics.Metrics{Service: mockServiceMetrics}}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(uint(5), nil)
	mockRepo.EXPECT().UpdateStockAmount(ctx,
		&dto.ArticleAmount{Article: "test-9", Amount: uint(4)}).Times(1).Return(nil)
	mockRepo.EXPECT().CreateReservation(ctx, &data).Times(1).Return(nil)
	mockServiceMetrics.EXPECT().PlacedInternetOrdersInc().Times(1)

	err := s.MakeReservation(ctx, data)
	if err != nil {
		t.Fail()
	}
}

func TestService_MakeReservationForLocalSuccess(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	mockServiceMetrics := mockService.NewMockMetricsInterface(ctrl)
	data := dto.NumberDateStateProducts{
		Products:    []dto.ArticlePriceAmount{{Article: "test-9", Amount: 1, Price: 698}},
		OrderNumber: reservation.MaxCashRegisterNumber + 1,
		Date:        time.Now(),
		State:       reservation.NewForLocalCustomer,
	}
	s := Service{Repository: mockRepo,
		Metrics: &metrics.Metrics{Service: mockServiceMetrics}}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(uint(5), nil)
	mockRepo.EXPECT().UpdateStockAmount(ctx,
		&dto.ArticleAmount{Article: "test-9", Amount: uint(4)}).Times(1).Return(nil)
	mockRepo.EXPECT().CreateReservation(ctx, &data).Times(1).Return(nil)
	mockServiceMetrics.EXPECT().PlacedLocalOrdersInc().Times(1)

	err := s.MakeReservation(ctx, data)
	if err != nil {
		t.Fail()
	}
}

func TestService_MakeReservationErrAmount(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.NumberDateStateProducts{
		Products:    []dto.ArticlePriceAmount{{Article: "test-9", Amount: 1, Price: 698}},
		OrderNumber: reservation.MaxCashRegisterNumber,
		Date:        time.Now(),
		State:       reservation.NewForCashRegister,
	}
	s := Service{Repository: mockRepo}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(
		uint(0), errors.New(""))

	err := s.MakeReservation(ctx, data)
	if err == nil {
		t.Fail()
	}
}

func TestService_MakeReservationNoEnough(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.NumberDateStateProducts{
		Products:    []dto.ArticlePriceAmount{{Article: "test-9", Amount: 2, Price: 698}},
		OrderNumber: reservation.MaxCashRegisterNumber,
		Date:        time.Now(),
		State:       reservation.NewForCashRegister,
	}
	s := Service{Repository: mockRepo}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(uint(1), nil)

	err := s.MakeReservation(ctx, data)
	if !errors.Is(err, service.ErrNoEnoughItemsToReserve) {
		t.Fail()
	}
}

func TestService_MakeReservationErrUpdate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.NumberDateStateProducts{
		Products:    []dto.ArticlePriceAmount{{Article: "test-9", Amount: 1, Price: 698}},
		OrderNumber: reservation.MaxCashRegisterNumber,
		Date:        time.Now(),
		State:       reservation.NewForCashRegister,
	}
	s := Service{Repository: mockRepo}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(uint(5), nil)
	mockRepo.EXPECT().UpdateStockAmount(ctx,
		&dto.ArticleAmount{Article: "test-9", Amount: uint(4)}).Times(1).Return(errors.New(""))

	err := s.MakeReservation(ctx, data)
	if err == nil {
		t.Fail()
	}
}

func TestService_MakeReservationErrCreate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.NumberDateStateProducts{
		Products:    []dto.ArticlePriceAmount{{Article: "test-9", Amount: 1, Price: 698}},
		OrderNumber: reservation.MaxCashRegisterNumber,
		Date:        time.Now(),
		State:       reservation.NewForCashRegister,
	}
	s := Service{Repository: mockRepo}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(uint(5), nil)
	mockRepo.EXPECT().UpdateStockAmount(ctx,
		&dto.ArticleAmount{Article: "test-9", Amount: uint(4)}).Times(1).Return(nil)
	mockRepo.EXPECT().CreateReservation(ctx, &data).Times(1).Return(errors.New(""))

	err := s.MakeReservation(ctx, data)
	if err == nil {
		t.Fail()
	}
}

func TestService_CancelReservationIncorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.Number{OrderNumber: 0}
	s := Service{Repository: mockRepo}

	mockRepo.EXPECT().WithinTransaction(context.Background(), gomock.Any()).Times(0)

	err := s.CancelReservation(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_CancelReservationCashRegister(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.Number{OrderNumber: 5}
	s := Service{Repository: mockRepo}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")

	mockRepo.EXPECT().ReadReservation(ctx, &data).Times(1).Return(
		dto.NumberDateStateProducts{
			Products:    []dto.ArticlePriceAmount{{Article: "test-9", Amount: 1, Price: 698}},
			OrderNumber: data.OrderNumber,
			Date:        time.Now(),
			State:       reservation.NewForCashRegister,
		}, nil)
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(uint(5), nil)
	mockRepo.EXPECT().UpdateStockAmount(ctx,
		&dto.ArticleAmount{Article: "test-9", Amount: uint(6)}).Times(1).Return(nil)
	mockRepo.EXPECT().DeleteReservation(ctx, &data).Times(1).Return(nil)

	err := s.CancelReservation(ctx, data)
	if err != nil {
		t.Fail()
	}
}

func TestService_CancelReservationErrReadReservation(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.Number{OrderNumber: 5}
	s := Service{Repository: mockRepo}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")

	mockRepo.EXPECT().ReadReservation(ctx, &data).Times(1).Return(dto.NumberDateStateProducts{}, repository.ErrNoRecord)

	err := s.CancelReservation(ctx, data)
	if !errors.Is(err, repository.ErrNoRecord) {
		t.Fail()
	}
}

func TestService_CancelReservationErrAlreadyFinished(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.Number{OrderNumber: 5}
	s := Service{Repository: mockRepo}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")

	mockRepo.EXPECT().ReadReservation(ctx, &data).Times(1).Return(
		dto.NumberDateStateProducts{
			Products:    []dto.ArticlePriceAmount{{Article: "test-9", Amount: 1, Price: 698}},
			OrderNumber: data.OrderNumber,
			Date:        time.Now(),
			State:       reservation.Finished,
		}, nil)

	err := s.CancelReservation(ctx, data)
	if !errors.Is(err, service.ErrAlreadyProcessed) {
		t.Fail()
	}
}

func TestService_CancelReservationErrReadStockAmount(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.Number{OrderNumber: 5}
	s := Service{Repository: mockRepo}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")

	mockRepo.EXPECT().ReadReservation(ctx, &data).Times(1).Return(
		dto.NumberDateStateProducts{
			Products:    []dto.ArticlePriceAmount{{Article: "test-9", Amount: 1, Price: 698}},
			OrderNumber: data.OrderNumber,
			Date:        time.Now(),
			State:       reservation.NewForCashRegister,
		}, nil)
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(
		uint(0), repository.ErrNoRecord)

	err := s.CancelReservation(ctx, data)
	if !errors.Is(err, repository.ErrNoRecord) {
		t.Fail()
	}
}

func TestService_CancelReservationErrUpdateAmount(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.Number{OrderNumber: 5}
	s := Service{Repository: mockRepo}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")

	mockRepo.EXPECT().ReadReservation(ctx, &data).Times(1).Return(
		dto.NumberDateStateProducts{
			Products:    []dto.ArticlePriceAmount{{Article: "test-9", Amount: 1, Price: 698}},
			OrderNumber: data.OrderNumber,
			Date:        time.Now(),
			State:       reservation.NewForCashRegister,
		}, nil)
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(uint(5), nil)
	mockRepo.EXPECT().UpdateStockAmount(ctx,
		&dto.ArticleAmount{Article: "test-9", Amount: uint(6)}).Times(1).Return(repository.ErrTimeout)

	err := s.CancelReservation(ctx, data)
	if !errors.Is(err, repository.ErrTimeout) {
		t.Fail()
	}
}

func TestService_CancelReservationInternet(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	mockServiceMetrics := mockService.NewMockMetricsInterface(ctrl)
	data := dto.Number{OrderNumber: 555}
	s := Service{Repository: mockRepo,
		Metrics: &metrics.Metrics{HTTP: nil, Service: mockServiceMetrics}}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	resData := dto.NumberDateStateProducts{Products: []dto.ArticlePriceAmount{{Article: "test-9", Amount: 1, Price: 698}},
		OrderNumber: 555, Date: time.Now(), State: reservation.NewForInternetCustomer,
	}
	mockRepo.EXPECT().ReadReservation(ctx, &data).Times(1).Return(resData, nil)
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(uint(5), nil)
	mockRepo.EXPECT().UpdateStockAmount(ctx,
		&dto.ArticleAmount{Article: "test-9", Amount: uint(6)}).Times(1).Return(nil)

	// Тест фейлился из-за расхождений во времени запуска time.Now() при создании DTO для функции UpdateReservation в
	// сервисе и тесте. Пришлось использовать в моке gomock.Any() вместо dto.NumberDateStateProducts
	mockRepo.EXPECT().UpdateReservation(ctx, gomock.Any()).Times(1).Return(nil)

	mockServiceMetrics.EXPECT().CancelOrdersInc().Times(1)

	err := s.CancelReservation(ctx, data)
	if err != nil {
		t.Fail()
	}
}

func TestService_MakeSaleErrDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := []dto.ArticlePriceAmount{{Article: "test-9.9999", Price: 410, Amount: 10}}
	s := Service{Repository: mockRepo}
	err := s.MakeSale(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_MakeSaleSuccess(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := []dto.ArticlePriceAmount{{Article: "test-9", Price: 410, Amount: 10}}
	s := Service{Repository: mockRepo}
	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")

	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(uint(12), nil)
	mockRepo.EXPECT().UpdateStockAmount(ctx, &dto.ArticleAmount{Article: "test-9", Amount: 2}).Times(1).Return(nil)
	mockRepo.EXPECT().CreateSoldRecord(ctx, gomock.Any()).Times(1).Return(nil)

	err := s.MakeSale(ctx, data)
	if err != nil {
		t.Fail()
	}
}

func TestService_MakeSaleErrCreateSoldRecord(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := []dto.ArticlePriceAmount{{Article: "test-9", Price: 410, Amount: 10}}
	s := Service{Repository: mockRepo}
	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")

	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(uint(12), nil)
	mockRepo.EXPECT().UpdateStockAmount(ctx, &dto.ArticleAmount{Article: "test-9", Amount: 2}).Times(1).Return(nil)
	mockRepo.EXPECT().CreateSoldRecord(ctx, gomock.Any()).Times(1).Return(repository.ErrTimeout)

	err := s.MakeSale(ctx, data)
	if !errors.Is(err, repository.ErrTimeout) {
		t.Fail()
	}
}

func TestService_MakeSaleErrUpdateStock(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := []dto.ArticlePriceAmount{{Article: "test-9", Price: 410, Amount: 10}}
	s := Service{Repository: mockRepo}
	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")

	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(uint(12), nil)
	mockRepo.EXPECT().UpdateStockAmount(ctx, &dto.ArticleAmount{Article: "test-9", Amount: 2}).Times(
		1).Return(repository.ErrTimeout)

	err := s.MakeSale(ctx, data)
	if !errors.Is(err, repository.ErrTimeout) {
		t.Fail()
	}
}

func TestService_MakeSaleErrNoEnoughItemsInStock(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := []dto.ArticlePriceAmount{{Article: "test-9", Price: 410, Amount: 10}}
	s := Service{Repository: mockRepo}
	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")

	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(uint(2), nil)

	err := s.MakeSale(ctx, data)
	if !errors.Is(err, service.ErrNoEnoughItemsInStock) {
		t.Fail()
	}
}

func TestService_MakeSaleErrReadStock(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := []dto.ArticlePriceAmount{{Article: "test-9", Price: 410, Amount: 10}}
	s := Service{Repository: mockRepo}
	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")

	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.Article{Article: "test-9"}).Times(1).Return(uint(2),
		repository.ErrTimeout)

	err := s.MakeSale(ctx, data)
	if !errors.Is(err, repository.ErrTimeout) {
		t.Fail()
	}
}

func TestService_FinishOrderIncorrectDTO(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	s := Service{Repository: mockRepo}

	err := s.FinishOrder(context.Background(), dto.Number{OrderNumber: 0})
	if err == nil {
		t.Fail()
	}
}

func TestService_FinishOrderLocal(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	s := Service{Repository: mockRepo}
	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	data := dto.Number{OrderNumber: reservation.MaxCashRegisterNumber}
	resData := dto.NumberDateStateProducts{
		Products:    []dto.ArticlePriceAmount{{Article: "test-9", Price: 100, Amount: 1}},
		OrderNumber: reservation.MaxCashRegisterNumber,
		Date:        time.Time{},
		State:       reservation.NewForCashRegister,
	}

	mockRepo.EXPECT().ReadReservation(ctx, &data).Times(1).Return(resData, nil)
	mockRepo.EXPECT().CreateSoldRecord(ctx, gomock.Any()).Times(1).Return(nil)
	mockRepo.EXPECT().DeleteReservation(ctx, &data).Times(1).Return(nil)

	err := s.FinishOrder(ctx, data)
	if err != nil {
		t.Fail()
	}
}

func TestService_FinishOrderInternet(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	s := Service{Repository: mockRepo}
	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	data := dto.Number{OrderNumber: reservation.MaxCashRegisterNumber + 1}
	resData := dto.NumberDateStateProducts{
		Products:    []dto.ArticlePriceAmount{{Article: "test-9", Price: 100, Amount: 1}},
		OrderNumber: reservation.MaxCashRegisterNumber + 1,
		Date:        time.Time{},
		State:       reservation.NewForInternetCustomer,
	}

	mockRepo.EXPECT().ReadReservation(ctx, &data).Times(1).Return(resData, nil)
	mockRepo.EXPECT().CreateSoldRecord(ctx, gomock.Any()).Times(1).Return(nil)
	mockRepo.EXPECT().UpdateReservation(ctx, gomock.Any()).Times(1).Return(nil)

	err := s.FinishOrder(ctx, data)
	if err != nil {
		t.Fail()
	}
}

func TestService_FinishOrderErrCreateSoldRecord(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	s := Service{Repository: mockRepo}
	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	data := dto.Number{OrderNumber: reservation.MaxCashRegisterNumber + 1}
	resData := dto.NumberDateStateProducts{
		Products:    []dto.ArticlePriceAmount{{Article: "test-9", Price: 100, Amount: 1}},
		OrderNumber: reservation.MaxCashRegisterNumber + 1,
		Date:        time.Time{},
		State:       reservation.NewForInternetCustomer,
	}

	mockRepo.EXPECT().ReadReservation(ctx, &data).Times(1).Return(resData, nil)
	mockRepo.EXPECT().CreateSoldRecord(ctx, gomock.Any()).Times(1).Return(repository.ErrTimeout)

	err := s.FinishOrder(ctx, data)
	if !errors.Is(err, repository.ErrTimeout) {
		t.Fail()
	}
}

func TestService_FinishOrderErrAlreadyProcessed(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	s := Service{Repository: mockRepo}
	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	data := dto.Number{OrderNumber: reservation.MaxCashRegisterNumber + 1}
	resData := dto.NumberDateStateProducts{
		Products:    []dto.ArticlePriceAmount{{Article: "test-9", Price: 100, Amount: 1}},
		OrderNumber: reservation.MaxCashRegisterNumber + 1,
		Date:        time.Time{},
		State:       reservation.Finished,
	}

	mockRepo.EXPECT().ReadReservation(ctx, &data).Times(1).Return(resData, nil)

	err := s.FinishOrder(ctx, data)
	if !errors.Is(err, service.ErrAlreadyProcessed) {
		t.Fail()
	}
}

func TestService_FinishOrderErrReadReservation(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	s := Service{Repository: mockRepo}
	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	data := dto.Number{OrderNumber: reservation.MaxCashRegisterNumber + 1}
	resData := dto.NumberDateStateProducts{
		Products:    []dto.ArticlePriceAmount{{Article: "test-9", Price: 100, Amount: 1}},
		OrderNumber: reservation.MaxCashRegisterNumber + 1,
		Date:        time.Time{},
		State:       reservation.NewForInternetCustomer,
	}

	mockRepo.EXPECT().ReadReservation(ctx, &data).Times(1).Return(resData, repository.ErrTimeout)

	err := s.FinishOrder(ctx, data)
	if !errors.Is(err, repository.ErrTimeout) {
		t.Fail()
	}
}
