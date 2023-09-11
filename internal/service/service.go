package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/lazylex/watch-store/store/internal/domain/aggregates/reservation"
	"github.com/lazylex/watch-store/store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store/store/internal/dto"
	"github.com/lazylex/watch-store/store/internal/helpers/constantes/prefixes"
	"github.com/lazylex/watch-store/store/internal/helpers/constantes/various"
	"github.com/lazylex/watch-store/store/internal/logger"
	"github.com/lazylex/watch-store/store/internal/ports/repository"
	standartLog "log"
	"log/slog"
	"time"
)

type Service struct {
	Repository repository.Interface
	Logger     *slog.Logger
}

// serviceError добавляет к тексту ошибки префикс, указывающий на её принадлежность к сервису
func serviceError(text string) error {
	return errors.New(prefixes.ServicePrefix + text)
}

var (
	ErrNoEnoughItemsToReserve = serviceError("no enough items to reserve")
	ErrNoEnoughItemsInStock   = serviceError("no enough items in stock")
	ErrAlreadyProcessed       = serviceError("already processed")
)

type Option func(*Service)

// WithLogger служит для подключения к сервису уже сконфигурированного логгера. Функция находится не в пакете logger,
// чтобы избежать перекрестных ссылок
func WithLogger(logger *slog.Logger) Option {
	return func(s *Service) {
		s.Logger = logger
	}
}

// New создаёт сервис. В качестве параметров передаются функции, инициализирующие в сервисе репозиторий с интерфейсом
// repository.Interface и логгер *slog.Logger
func New(options ...Option) *Service {
	requiredOptions, initializedOptions := 2, 0

	service := &Service{}
	for _, opt := range options {
		opt(service)
		initializedOptions++
	}

	if initializedOptions != requiredOptions {
		standartLog.Fatal(serviceError(
			fmt.Sprintf("need to initialize %d options, not %d", requiredOptions, initializedOptions)).Error())
	}

	return service
}

// ChangePriceInStock изменяет цену товара, находящегося в продаже
func (s *Service) ChangePriceInStock(ctx context.Context, data dto.ArticleWithPriceDTO) error {
	if err := data.Validate(); err != nil {
		return err
	}

	err := s.Repository.UpdateStockPrice(ctx, &data)
	if err == nil {
		logger.LogWithCtxData(ctx, s.Logger.With(logger.OPLabel, "service.ChangePriceInStock")).Info(
			fmt.Sprintf("change price to %.2f in stock record with article %s", data.Price, data.Article))
	}
	return err
}

// GetStock возвращает полную информацию о товаре, доступном для продажи, в виде dto.NamedProductDTO
func (s *Service) GetStock(ctx context.Context, data dto.ArticleDTO) (dto.NamedProductDTO, error) {
	if err := data.Validate(); err != nil {
		return dto.NamedProductDTO{}, err
	}
	sale, err := s.Repository.ReadStock(ctx, &data)
	if err != nil {
		return dto.NamedProductDTO{}, err
	}

	logger.LogWithCtxData(ctx, s.Logger.With(logger.OPLabel, "service.GetStock")).Info(
		fmt.Sprintf("requested stock record with article %s", data.Article))

	return sale, nil
}

// AddProductToStock добавляет новый товар в ассортимент магазина
func (s *Service) AddProductToStock(ctx context.Context, data dto.NamedProductDTO) error {
	if err := data.Validate(); err != nil {
		return err
	}

	if err := s.Repository.CreateStock(ctx, &data); err != nil {
		return err
	}

	logger.LogWithCtxData(ctx, s.Logger.With(logger.OPLabel, "service.AddProductToStock")).Info(
		fmt.Sprintf("add to stock record with article %s, price %.2f", data.Article, data.Price))
	return nil
}

// ChangeAmountInStock изменяет доступное для продажи количество товара
func (s *Service) ChangeAmountInStock(ctx context.Context, data dto.ArticleWithAmountDTO) error {
	if err := data.Validate(); err != nil {
		return err
	}

	if err := s.Repository.UpdateStockAmount(ctx, &data); err != nil {
		return err
	}

	logger.LogWithCtxData(ctx, s.Logger.With(logger.OPLabel, "service.ChangeAmountInStock")).Info(
		fmt.Sprintf("amount udpaded to %d in stock record with article %s", data.Amount, data.Article))
	return nil
}

// GetAmountInStock возвращает доступное для продажи количество товара
func (s *Service) GetAmountInStock(ctx context.Context, data dto.ArticleDTO) (uint, error) {
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
func (s *Service) MakeReservation(ctx context.Context, data dto.ReservationDTO) error {
	var err error
	var available uint
	newAmountInStock := make(map[article.Article]uint)

	if err = data.Validate(); err != nil {
		return err
	}

	return s.Repository.WithinTransaction(ctx, func(txCtx context.Context) error {
		for _, p := range data.Products {
			if available, err = s.Repository.ReadStockAmount(txCtx,
				&dto.ArticleDTO{Article: p.Article}); err != nil {
				return err
			}
			if available < p.Amount {
				return ErrNoEnoughItemsToReserve
			}
			newAmountInStock[p.Article] = available - p.Amount
		}
		for _, p := range data.Products {
			err = s.Repository.UpdateStockAmount(
				txCtx,
				&dto.ArticleWithAmountDTO{
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

		logger.LogWithCtxData(txCtx, s.Logger.With(logger.OPLabel, "service.MakeReservation")).Info(
			fmt.Sprintf("succesfully saved order %d", data.OrderNumber))
		return nil
	})
}

// CancelReservation снимает бронь с товара/ов
func (s *Service) CancelReservation(ctx context.Context, data dto.OrderNumberDTO) error {
	if err := data.Validate(); err != nil {
		return err
	}

	return s.Repository.WithinTransaction(ctx, func(txCtx context.Context) error {
		res, err := s.Repository.ReadReservation(txCtx, &data)
		if err != nil {
			return err
		}

		if !res.IsNew() {
			return ErrAlreadyProcessed
		}

		for _, p := range res.Products {
			var inStock uint
			if inStock, err = s.Repository.ReadStockAmount(txCtx, &dto.ArticleDTO{Article: p.Article}); err != nil {
				return err
			}
			if err = s.Repository.UpdateStockAmount(txCtx,
				&dto.ArticleWithAmountDTO{Article: p.Article, Amount: p.Amount + inStock}); err != nil {
				return err
			}

		}

		if data.OrderNumber <= reservation.MaxCashRegisterNumber {
			return s.Repository.DeleteReservation(txCtx, &dto.OrderNumberDTO{OrderNumber: data.OrderNumber})
		}

		return s.Repository.UpdateReservation(txCtx, &dto.ReservationDTO{
			Products:    res.Products,
			OrderNumber: data.OrderNumber,
			Date:        time.Now(),
			State:       reservation.Cancel,
		})
	})
}

// MakeSale уменьшает количества доступного для продажи товара и производит запись в статистику продаж
func (s *Service) MakeSale(ctx context.Context, data []dto.ProductDTO) error {
	for _, p := range data {
		if err := p.Validate(); err != nil {
			return err
		}
	}

	var err error
	var available uint

	return s.Repository.WithinTransaction(ctx, func(txCtx context.Context) error {
		for _, p := range data {
			if available, err = s.Repository.ReadStockAmount(txCtx, &dto.ArticleDTO{Article: p.Article}); err != nil {
				return err
			}
			if available < p.Amount {
				return ErrNoEnoughItemsInStock
			}

			if err = s.Repository.UpdateStockAmount(txCtx, &dto.ArticleWithAmountDTO{
				Article: p.Article,
				Amount:  available - p.Amount,
			},
			); err != nil {
				return err
			}
		}

		for _, p := range data {
			if err = s.Repository.CreateSoldRecord(txCtx, &dto.SoldDTO{
				Article: p.Article, Price: p.Price, Amount: p.Amount, Date: time.Now(),
			}); err != nil {
				return err
			}
		}

		logger.LogWithCtxData(txCtx, s.Logger.With(logger.OPLabel, "service.MakeSale")).Info(
			"sale completed successfully")

		return nil
	})
}

// FinishOrder помечает заказ, как выполненный. Данные о содержащихся в заказе товарах переносятся в статистику продаж
func (s *Service) FinishOrder(ctx context.Context, data dto.OrderNumberDTO) error {
	if err := data.Validate(); err != nil {
		return err
	}

	return s.Repository.WithinTransaction(ctx, func(txCtx context.Context) error {

		res, err := s.Repository.ReadReservation(txCtx, &data)
		if err != nil {
			return err
		}

		if !res.IsNew() {
			return ErrAlreadyProcessed
		}

		for _, p := range res.Products {
			if err = s.Repository.CreateSoldRecord(txCtx, &dto.SoldDTO{
				Article: p.Article,
				Price:   p.Price,
				Amount:  p.Amount,
				Date:    time.Now(),
			}); err != nil {
				return err
			}
		}

		if data.OrderNumber <= reservation.MaxCashRegisterNumber {
			return s.Repository.DeleteReservation(txCtx, &dto.OrderNumberDTO{OrderNumber: data.OrderNumber})
		}

		return s.Repository.UpdateReservation(txCtx, &dto.ReservationDTO{
			Products:    res.Products,
			OrderNumber: data.OrderNumber,
			Date:        time.Now(),
			State:       reservation.Finished,
		})
	})
}

// TotalSold возвращает количество проданного товара с переданным артикулом за весь период
func (s *Service) TotalSold(ctx context.Context, data dto.ArticleDTO) (uint, error) {
	var amount uint
	var err error

	if err = data.Validate(); err != nil {
		return 0, err
	}

	if amount, err = s.Repository.ReadSoldAmount(ctx, &data); err != nil {
		return 0, err
	}

	logger.LogWithCtxData(ctx, s.Logger.With(logger.OPLabel, "service.TotalSold")).Info(
		fmt.Sprintf("readed amount of sold %d (article %s)", amount, data.Article))

	return amount, nil
}

// TotalSoldInPeriod возвращает количество проданного товара с переданным артикулом за указанный период
func (s *Service) TotalSoldInPeriod(ctx context.Context, data dto.ArticleWithPeriodDTO) (uint, error) {
	var amount uint
	var err error

	if err = data.Validate(); err != nil {
		return 0, err
	}

	if amount, err = s.Repository.ReadSoldAmountInPeriod(ctx, &data); err != nil {
		return 0, err
	}

	logger.LogWithCtxData(ctx, s.Logger.With(logger.OPLabel, "service.TotalSoldInPeriod")).Info(
		fmt.Sprintf("readed amount of sold - %d (from %s to %s) (article %s)",
			amount, data.From.Format(various.DateLayout), data.To.Format(various.DateLayout), data.Article))

	return amount, nil
}
