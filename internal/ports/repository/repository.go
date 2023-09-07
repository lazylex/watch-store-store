package repository

import (
	"context"
	"errors"
	"github.com/lazylex/watch-store/store/internal/dto"
	"github.com/lazylex/watch-store/store/internal/helpers/constantes/prefixes"
)

// repositoryError добавляет к тексту ошибки префикс, указывающий на её принадлежность к хранилищу
func repositoryError(text string) error {
	return errors.New(prefixes.RepositoryPrefix + text)
}

var (
	ErrNoRecord = repositoryError("no record")
	ErrTimeout  = repositoryError("operation timeout")
)

type Interface interface {
	ConvertToCommonErr(error) error
	Close() error

	WithinTransaction(context.Context, func(context.Context) error) error

	CreateStock(context.Context, *dto.NamedProductDTO) error
	ReadStock(context.Context, *dto.ArticleDTO) (dto.NamedProductDTO, error)
	ReadStockAmount(context.Context, *dto.ArticleDTO) (uint, error)
	ReadStockPrice(context.Context, *dto.ArticleDTO) (float64, error)
	UpdateStock(context.Context, *dto.NamedProductDTO) error
	UpdateStockAmount(context.Context, *dto.ArticleWithAmountDTO) error
	UpdateStockPrice(context.Context, *dto.ArticleWithPriceDTO) error

	CreateReservation(context.Context, *dto.ReservationDTO) error
	ReadReservation(context.Context, *dto.OrderNumberDTO) (dto.ReservationDTO, error)
	UpdateReservation(context.Context, *dto.ReservationDTO) error
	DeleteReservation(context.Context, *dto.OrderNumberDTO) error

	CreateSoldRecord(context.Context, *dto.SoldDTO) error
	ReadSoldRecords(context.Context, *dto.ArticleDTO) ([]dto.SoldDTO, error)
	ReadSoldRecordsInPeriod(context.Context, *dto.ArticleWithPeriodDTO) ([]dto.SoldDTO, error)
}
