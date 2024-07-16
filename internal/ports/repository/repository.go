package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lazylex/watch-store-store/internal/dto"
	"github.com/lazylex/watch-store-store/internal/helpers/constants/prefixes"
)

// repositoryError добавляет к тексту ошибки префикс, указывающий на её принадлежность к хранилищу.
func repositoryError(text string) error {
	return errors.New(prefixes.RepositoryPrefix + text)
}

var (
	ErrNoRecord  = repositoryError("no record")
	ErrTimeout   = repositoryError("operation timeout")
	ErrDuplicate = repositoryError("duplicate entry")
)

//go:generate mockgen -source=repository.go -destination=mocks/repository.go
type Interface interface {
	ConvertToCommonErr(error) error

	WithinTransaction(context.Context, func(context.Context) error) error

	CreateStock(context.Context, *dto.ArticlePriceNameAmount) error
	ReadStock(context.Context, *dto.Article) (dto.ArticlePriceNameAmount, error)
	ReadStockAmount(context.Context, *dto.Article) (uint, error)
	ReadStockPrice(context.Context, *dto.Article) (float64, error)
	UpdateStock(context.Context, *dto.ArticlePriceNameAmount) error
	UpdateStockAmount(context.Context, *dto.ArticleAmount) error
	UpdateStockPrice(context.Context, *dto.ArticlePrice) error

	CreateReservation(context.Context, *dto.NumberDateStateProducts) error
	ReadReservation(context.Context, *dto.Number) (dto.NumberDateStateProducts, error)
	UpdateReservation(context.Context, *dto.NumberDateStateProducts) error
	DeleteReservation(context.Context, *dto.Number) error

	CreateSoldRecord(context.Context, *dto.ArticlePriceAmountDate) error
	ReadSoldRecords(context.Context, *dto.Article) ([]dto.ArticlePriceAmountDate, error)
	ReadSoldAmount(context.Context, *dto.Article) (uint, error)
	ReadSoldRecordsInPeriod(context.Context, *dto.ArticleFromTo) ([]dto.ArticlePriceAmountDate, error)
	ReadSoldAmountInPeriod(context.Context, *dto.ArticleFromTo) (uint, error)
}

type SQLDBInterface interface {
	DB() *sql.DB
	Close() error
}
