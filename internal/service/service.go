package service

import (
	"context"
	"fmt"
	"github.com/lazylex/watch-store-store/internal/domain/aggregates/reservation"
	"github.com/lazylex/watch-store-store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store-store/internal/dto"
	"github.com/lazylex/watch-store-store/internal/helpers/constants/prefixes"
	"github.com/lazylex/watch-store-store/internal/helpers/constants/various"
	"github.com/lazylex/watch-store-store/internal/logger"
	"github.com/lazylex/watch-store-store/internal/metrics"
	"github.com/lazylex/watch-store-store/internal/ports/repository"
	"github.com/lazylex/watch-store-store/internal/ports/service"
	standartLog "log"
	"log/slog"
	"time"
)

type Service struct {
	Repository    repository.Interface
	SQLRepository repository.SQLDBInterface
	Metrics       *metrics.Metrics
}

type Option func(*Service)

// WithMetrics служит для внедрения в сервис уже инициализированных метрик для Prometheus.
func WithMetrics(metrics *metrics.Metrics) Option {
	return func(s *Service) {
		s.Metrics = metrics
	}
}

// New создаёт сервис. В качестве параметров передаются функции, инициализирующие в сервисе репозиторий с интерфейсом
// repository.Interface и структура работы с логами *slog.Logger.
func New(options ...Option) *Service {
	requiredOptions, initializedOptions := 2, 0

	s := &Service{}
	for _, opt := range options {
		opt(s)
		initializedOptions++
	}

	if initializedOptions != requiredOptions {
		standartLog.Fatal(prefixes.ServicePrefix +
			fmt.Sprintf("need to initialize %d options, not %d", requiredOptions, initializedOptions))
	}

	return s
}

// ChangePriceInStock изменяет цену товара, находящегося в продаже.
func (s *Service) ChangePriceInStock(ctx context.Context, data dto.ArticlePrice) error {
	if err := data.Validate(); err != nil {
		return err
	}
	_, err := s.Stock(ctx, dto.Article{Article: data.Article})
	if err != nil {
		return err
	}

	err = s.Repository.UpdateStockPrice(ctx, &data)
	if err == nil {
		logger.LogWithCtxData(ctx, slog.With(logger.OPLabel, "service.ChangePriceInStock")).Info(
			fmt.Sprintf("change price to %.2f in stock record with article %s", data.Price, data.Article))
	}
	return err
}

// Stock возвращает полную информацию о товаре, доступном для продажи, в виде dto.ArticlePriceNameAmount.
func (s *Service) Stock(ctx context.Context, data dto.Article) (dto.ArticlePriceNameAmount, error) {
	if err := data.Validate(); err != nil {
		return dto.ArticlePriceNameAmount{}, err
	}
	sale, err := s.Repository.ReadStock(ctx, &data)
	if err != nil {
		return dto.ArticlePriceNameAmount{}, err
	}

	logger.LogWithCtxData(ctx, slog.With(logger.OPLabel, "service.Stock")).Info(
		fmt.Sprintf("requested stock record with article %s", data.Article))

	return sale, nil
}

// AddProductToStock добавляет новый товар в ассортимент магазина.
func (s *Service) AddProductToStock(ctx context.Context, data dto.ArticlePriceNameAmount) error {
	if err := data.Validate(); err != nil {
		return err
	}

	if err := s.Repository.CreateStock(ctx, &data); err != nil {
		return err
	}

	logger.LogWithCtxData(ctx, slog.With(logger.OPLabel, "service.AddProductToStock")).Info(
		fmt.Sprintf("add to stock record with article %s, price %.2f", data.Article, data.Price))
	return nil
}

// ChangeAmountInStock изменяет доступное для продажи количество товара.
func (s *Service) ChangeAmountInStock(ctx context.Context, data dto.ArticleAmount) error {
	if err := data.Validate(); err != nil {
		return err
	}

	if err := s.Repository.UpdateStockAmount(ctx, &data); err != nil {
		return err
	}

	logger.LogWithCtxData(ctx, slog.With(logger.OPLabel, "service.ChangeAmountInStock")).Info(
		fmt.Sprintf("amount udpaded to %d in stock record with article %s", data.Amount, data.Article))
	return nil
}

// AmountInStock возвращает доступное для продажи количество товара.
func (s *Service) AmountInStock(ctx context.Context, data dto.Article) (uint, error) {
	if err := data.Validate(); err != nil {
		return 0, err
	}

	amount, err := s.Repository.ReadStockAmount(ctx, &data)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

// MakeReservation производит резервирование товара для покупателя. Резервирование проводится как для бронирования
// через интернет, так и во время нахождения товара на кассе (в ожидании оплаты локальным покупателем). В таком случае
// в качестве номера заказа передаётся номер кассы.
func (s *Service) MakeReservation(ctx context.Context, data dto.NumberDateStateProducts) error {
	var err error
	var available uint
	newAmountInStock := make(map[article.Article]uint)

	if err = data.Validate(); err != nil {
		return err
	}

	return s.Repository.WithinTransaction(ctx, func(txCtx context.Context) error {
		for _, p := range data.Products {
			if available, err = s.Repository.ReadStockAmount(txCtx,
				&dto.Article{Article: p.Article}); err != nil {
				return err
			}
			if available < p.Amount {
				return service.ErrNoEnoughItemsToReserve
			}
			newAmountInStock[p.Article] = available - p.Amount
		}
		for _, p := range data.Products {
			err = s.Repository.UpdateStockAmount(
				txCtx,
				&dto.ArticleAmount{
					Article: p.Article,
					Amount:  newAmountInStock[p.Article],
				},
			)
			if err != nil {
				return err
			}
		}
		if err = s.Repository.CreateReservation(txCtx, &data); err != nil {
			return err
		}

		if data.State == reservation.NewForInternetCustomer {
			s.Metrics.Service.PlacedInternetOrdersInc()
		}

		if data.State == reservation.NewForLocalCustomer {
			s.Metrics.Service.PlacedLocalOrdersInc()
		}

		logger.LogWithCtxData(txCtx, slog.With(logger.OPLabel, "service.MakeReservation")).Info(
			fmt.Sprintf("succesfully saved order %d", data.OrderNumber))
		return nil
	})
}

// CancelReservation снимает бронь с товара/ов.
func (s *Service) CancelReservation(ctx context.Context, data dto.Number) error {
	if err := data.Validate(); err != nil {
		return err
	}

	return s.Repository.WithinTransaction(ctx, func(txCtx context.Context) error {
		res, err := s.Repository.ReadReservation(txCtx, &data)
		if err != nil {
			return err
		}

		if !res.IsNew() {
			return service.ErrAlreadyProcessed
		}

		for _, p := range res.Products {
			var inStock uint
			if inStock, err = s.Repository.ReadStockAmount(txCtx, &dto.Article{Article: p.Article}); err != nil {
				return err
			}
			if err = s.Repository.UpdateStockAmount(txCtx,
				&dto.ArticleAmount{Article: p.Article, Amount: p.Amount + inStock}); err != nil {
				return err
			}

		}

		if data.OrderNumber <= reservation.MaxCashRegisterNumber {
			return s.Repository.DeleteReservation(txCtx, &dto.Number{OrderNumber: data.OrderNumber})
		}

		err = s.Repository.UpdateReservation(txCtx, &dto.NumberDateStateProducts{
			Products:    res.Products,
			OrderNumber: data.OrderNumber,
			Date:        time.Now(),
			State:       reservation.Cancel,
		})

		if err == nil {
			s.Metrics.Service.CancelOrdersInc()
		}

		return err
	})
}

// MakeSale уменьшает количества доступного для продажи товара и производит запись в статистику продаж.
func (s *Service) MakeSale(ctx context.Context, data []dto.ArticlePriceAmount) error {
	for _, p := range data {
		if err := p.Validate(); err != nil {
			return err
		}
	}

	var err error
	var available uint

	return s.Repository.WithinTransaction(ctx, func(txCtx context.Context) error {
		for _, p := range data {
			if available, err = s.Repository.ReadStockAmount(txCtx, &dto.Article{Article: p.Article}); err != nil {
				return err
			}
			if available < p.Amount {
				return service.ErrNoEnoughItemsInStock
			}

			if err = s.Repository.UpdateStockAmount(txCtx, &dto.ArticleAmount{
				Article: p.Article,
				Amount:  available - p.Amount,
			},
			); err != nil {
				return err
			}
		}

		for _, p := range data {
			if err = s.Repository.CreateSoldRecord(txCtx, &dto.ArticlePriceAmountDate{
				Article: p.Article, Price: p.Price, Amount: p.Amount, Date: time.Now(),
			}); err != nil {
				return err
			}
		}

		logger.LogWithCtxData(txCtx, slog.With(logger.OPLabel, "service.MakeSale")).Info(
			"sale completed successfully")

		return nil
	})
}

// FinishOrder помечает заказ, как выполненный. Данные о содержащихся в заказе товарах переносятся в статистику продаж.
func (s *Service) FinishOrder(ctx context.Context, data dto.Number) error {
	if err := data.Validate(); err != nil {
		return err
	}

	return s.Repository.WithinTransaction(ctx, func(txCtx context.Context) error {

		res, err := s.Repository.ReadReservation(txCtx, &data)
		if err != nil {
			return err
		}

		if !res.IsNew() {
			return service.ErrAlreadyProcessed
		}

		for _, p := range res.Products {
			if err = s.Repository.CreateSoldRecord(txCtx, &dto.ArticlePriceAmountDate{
				Article: p.Article,
				Price:   p.Price,
				Amount:  p.Amount,
				Date:    time.Now(),
			}); err != nil {
				return err
			}
		}

		if data.OrderNumber <= reservation.MaxCashRegisterNumber {
			return s.Repository.DeleteReservation(txCtx, &data)
		}

		return s.Repository.UpdateReservation(txCtx, &dto.NumberDateStateProducts{
			Products:    res.Products,
			OrderNumber: data.OrderNumber,
			Date:        time.Now(),
			State:       reservation.Finished,
		})
	})
}

// TotalSold возвращает количество проданного товара с переданным артикулом за весь период.
func (s *Service) TotalSold(ctx context.Context, data dto.Article) (uint, error) {
	var amount uint
	var err error

	if err = data.Validate(); err != nil {
		return 0, err
	}

	if amount, err = s.Repository.ReadSoldAmount(ctx, &data); err != nil {
		return 0, err
	}

	logger.LogWithCtxData(ctx, slog.With(logger.OPLabel, "service.TotalSold")).Info(
		fmt.Sprintf("readed amount of sold %d (article %s)", amount, data.Article))

	return amount, nil
}

// TotalSoldInPeriod возвращает количество проданного товара с переданным артикулом за указанный период.
func (s *Service) TotalSoldInPeriod(ctx context.Context, data dto.ArticleFromTo) (uint, error) {
	var amount uint
	var err error

	if err = data.Validate(); err != nil {
		return 0, err
	}

	if amount, err = s.Repository.ReadSoldAmountInPeriod(ctx, &data); err != nil {
		return 0, err
	}

	logger.LogWithCtxData(ctx, slog.With(logger.OPLabel, "service.TotalSoldInPeriod")).Info(
		fmt.Sprintf("readed amount of sold - %d (from %s to %s) (article %s)",
			amount, data.From.Format(various.DateLayout), data.To.Format(various.DateLayout), data.Article))

	return amount, nil
}
