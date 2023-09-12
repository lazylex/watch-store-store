package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lazylex/watch-store/store/internal/domain/aggregates/reservation"
	"github.com/lazylex/watch-store/store/internal/dto"
	"github.com/lazylex/watch-store/store/internal/ports/repository"
	mockrepository "github.com/lazylex/watch-store/store/internal/ports/repository/mocks"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"testing"
)

var nullLogger = slog.New(slog.NewTextHandler(os.NewFile(0, os.DevNull), nil))

// withMockRepo позволяет подключать мок репозитория в функции service.New, чтобы покрыть её тестами
func withMockRepo(mr repository.Interface) Option {
	return func(s *Service) {
		s.Repository = mr
	}
}

func TestNew(t *testing.T) {
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
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.NamedProductDTO{Name: "test_correct", Article: "test-9", Price: 110, Amount: 10}
	s := New(WithLogger(nullLogger), withMockRepo(mockRepo))

	mockRepo.EXPECT().CreateStock(context.Background(), &data).Times(1).Return(nil)

	err := s.AddProductToStock(context.Background(), data)
	if err != nil {
		t.Fail()
	}
}

func TestService_AddProductToStockIncorrectDTO(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.NamedProductDTO{Name: "test_incorrect", Article: "test-9.9999", Price: -110, Amount: 10}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().CreateStock(context.Background(), &data).Times(0)

	err := s.AddProductToStock(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_AddProductToStockDuplicateArticle(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.NamedProductDTO{Name: "test_correct", Article: "test-9", Price: 110, Amount: 10}
	s := Service{Repository: mockRepo, Logger: nullLogger}

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
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleWithPriceDTO{Article: "test-9", Price: 10}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().UpdateStockPrice(context.Background(), &data).Times(1).Return(nil)

	err := s.ChangePriceInStock(context.Background(), data)
	if err != nil {
		t.Fail()
	}
}

func TestService_ChangePriceInStockIncorrectDTO(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleWithPriceDTO{Article: "test-9", Price: -10}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().UpdateStockPrice(context.Background(), &data).Times(0)

	err := s.ChangePriceInStock(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_GetStockCorrectDTO(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleDTO{Article: "test-9"}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().ReadStock(context.Background(), &data).Times(1).Return(dto.NamedProductDTO{
		Name: "test-9", Article: "test-9", Price: 110, Amount: 10}, nil)

	_, err := s.GetStock(context.Background(), data)
	if err != nil {
		t.Fail()
	}
}

func TestService_GetStockIncorrectDTO(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleDTO{Article: "test-9.9999"}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().ReadStock(context.Background(), &data).Times(0)

	_, err := s.GetStock(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_GetStockNoRecord(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleDTO{Article: "test-9"}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().ReadStock(context.Background(), &data).Times(1).Return(dto.NamedProductDTO{},
		repository.ErrNoRecord)

	_, err := s.GetStock(context.Background(), data)
	if !errors.Is(err, repository.ErrNoRecord) {
		t.Fail()
	}
}

func TestService_ChangeAmountInStockCorrectDTO(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleWithAmountDTO{Article: "test-9", Amount: 10}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().UpdateStockAmount(context.Background(), &data).Times(1).Return(nil)

	err := s.ChangeAmountInStock(context.Background(), data)
	if err != nil {
		t.Fail()
	}
}

func TestService_ChangeAmountInStockIncorrectDTO(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleWithAmountDTO{Article: "test-9.9999", Amount: 10}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().UpdateStockAmount(context.Background(), &data).Times(0)

	err := s.ChangeAmountInStock(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_ChangeAmountInStockNoRecord(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleWithAmountDTO{Article: "test-9", Amount: 10}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().UpdateStockAmount(context.Background(), &data).Times(1).Return(repository.ErrNoRecord)

	err := s.ChangeAmountInStock(context.Background(), data)
	if !errors.Is(err, repository.ErrNoRecord) {
		t.Fail()
	}
}

func TestService_GetAmountInStockCorrectDTO(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleDTO{Article: "test-9"}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().ReadStockAmount(context.Background(), &data).Times(1).Return(uint(5), nil)

	_, err := s.GetAmountInStock(context.Background(), data)
	if err != nil {
		t.Fail()
	}
}

func TestService_GetAmountInStockIncorrectDTO(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleDTO{Article: "test-9.9999"}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().ReadStockAmount(context.Background(), &data).Times(0)

	_, err := s.GetAmountInStock(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_GetAmountInStockNoRecord(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleDTO{Article: "test-9"}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().ReadStockAmount(context.Background(), &data).Times(1).Return(uint(0), repository.ErrNoRecord)

	_, err := s.GetAmountInStock(context.Background(), data)
	if !errors.Is(err, repository.ErrNoRecord) {
		t.Fail()
	}
}

func TestService_TotalSoldCorrectDTO(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleDTO{Article: "test-9"}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().ReadSoldAmount(context.Background(), &data).Times(1).Return(uint(5), nil)

	_, err := s.TotalSold(context.Background(), data)
	if err != nil {
		t.Fail()
	}
}

func TestService_TotalSoldIncorrectDTO(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleDTO{Article: "test-9.9999"}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().ReadSoldAmount(context.Background(), &data).Times(0)

	_, err := s.TotalSold(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_TotalSoldTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleDTO{Article: "test-9"}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().ReadSoldAmount(context.Background(), &data).Times(1).Return(uint(0), repository.ErrTimeout)

	_, err := s.TotalSold(context.Background(), data)
	if !errors.Is(err, repository.ErrTimeout) {
		t.Fail()
	}
}

func TestService_TotalSoldInPeriodCorrectDTO(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleWithPeriodDTO{
		Article: "test_9",
		From:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		To:      time.Now(),
	}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().ReadSoldAmountInPeriod(context.Background(), &data).Times(1).Return(uint(5), nil)

	_, err := s.TotalSoldInPeriod(context.Background(), data)
	if err != nil {
		t.Fail()
	}
}

func TestService_TotalSoldInPeriodIncorrectDTO(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleWithPeriodDTO{Article: "test_9.9999", From: time.Now(), To: time.Now()}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().ReadSoldAmountInPeriod(context.Background(), &data).Times(0)

	_, err := s.TotalSoldInPeriod(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_TotalSoldInPeriodTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ArticleWithPeriodDTO{
		Article: "test_9",
		From:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		To:      time.Now(),
	}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().ReadSoldAmountInPeriod(context.Background(), &data).Times(1).Return(
		uint(0), repository.ErrTimeout)

	_, err := s.TotalSoldInPeriod(context.Background(), data)
	if !errors.Is(err, repository.ErrTimeout) {
		t.Fail()
	}
}

func TestService_MakeReservationIncorrectDTO(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ReservationDTO{
		Products:    []dto.ProductDTO{{Article: "test-9", Amount: 1, Price: 698}},
		OrderNumber: reservation.MaxCashRegisterNumber + 1,
		Date:        time.Now(),
		State:       reservation.NewForCashRegister,
	}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	mockRepo.EXPECT().WithinTransaction(context.Background(), gomock.Any()).Times(0)

	err := s.MakeReservation(context.Background(), data)
	if err == nil {
		t.Fail()
	}
}

func TestService_MakeReservationSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ReservationDTO{
		Products:    []dto.ProductDTO{{Article: "test-9", Amount: 1, Price: 698}},
		OrderNumber: reservation.MaxCashRegisterNumber,
		Date:        time.Now(),
		State:       reservation.NewForCashRegister,
	}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.ArticleDTO{Article: "test-9"}).Times(1).Return(uint(5), nil)
	mockRepo.EXPECT().UpdateStockAmount(ctx,
		&dto.ArticleWithAmountDTO{Article: "test-9", Amount: uint(4)}).Times(1).Return(nil)
	mockRepo.EXPECT().CreateReservation(ctx, &data).Times(1).Return(nil)

	err := s.MakeReservation(ctx, data)
	if err != nil {
		t.Fail()
	}
}

func TestService_MakeReservationErrAmount(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ReservationDTO{
		Products:    []dto.ProductDTO{{Article: "test-9", Amount: 1, Price: 698}},
		OrderNumber: reservation.MaxCashRegisterNumber,
		Date:        time.Now(),
		State:       reservation.NewForCashRegister,
	}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.ArticleDTO{Article: "test-9"}).Times(1).Return(
		uint(0), errors.New(""))

	err := s.MakeReservation(ctx, data)
	if err == nil {
		t.Fail()
	}
}

func TestService_MakeReservationNoEnough(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ReservationDTO{
		Products:    []dto.ProductDTO{{Article: "test-9", Amount: 2, Price: 698}},
		OrderNumber: reservation.MaxCashRegisterNumber,
		Date:        time.Now(),
		State:       reservation.NewForCashRegister,
	}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.ArticleDTO{Article: "test-9"}).Times(1).Return(uint(1), nil)

	err := s.MakeReservation(ctx, data)
	if !errors.Is(err, ErrNoEnoughItemsToReserve) {
		t.Fail()
	}
}

func TestService_MakeReservationErrUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ReservationDTO{
		Products:    []dto.ProductDTO{{Article: "test-9", Amount: 1, Price: 698}},
		OrderNumber: reservation.MaxCashRegisterNumber,
		Date:        time.Now(),
		State:       reservation.NewForCashRegister,
	}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.ArticleDTO{Article: "test-9"}).Times(1).Return(uint(5), nil)
	mockRepo.EXPECT().UpdateStockAmount(ctx,
		&dto.ArticleWithAmountDTO{Article: "test-9", Amount: uint(4)}).Times(1).Return(errors.New(""))

	err := s.MakeReservation(ctx, data)
	if err == nil {
		t.Fail()
	}
}

func TestService_MakeReservationErrCreate(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockrepository.NewMockInterface(ctrl)
	data := dto.ReservationDTO{
		Products:    []dto.ProductDTO{{Article: "test-9", Amount: 1, Price: 698}},
		OrderNumber: reservation.MaxCashRegisterNumber,
		Date:        time.Now(),
		State:       reservation.NewForCashRegister,
	}
	s := Service{Repository: mockRepo, Logger: nullLogger}

	ctx := context.WithValue(context.Background(), mockrepository.ExecuteKey{}, "✅")
	mockRepo.EXPECT().ReadStockAmount(ctx, &dto.ArticleDTO{Article: "test-9"}).Times(1).Return(uint(5), nil)
	mockRepo.EXPECT().UpdateStockAmount(ctx,
		&dto.ArticleWithAmountDTO{Article: "test-9", Amount: uint(4)}).Times(1).Return(nil)
	mockRepo.EXPECT().CreateReservation(ctx, &data).Times(1).Return(errors.New(""))

	err := s.MakeReservation(ctx, data)
	if err == nil {
		t.Fail()
	}
}
