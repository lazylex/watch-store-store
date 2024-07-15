package service

import (
	"context"
	"errors"
	"github.com/lazylex/watch-store-store/internal/dto"
	"github.com/lazylex/watch-store-store/internal/helpers/constants/prefixes"
)

// serviceError добавляет к тексту ошибки префикс, указывающий на её принадлежность к сервису.
func serviceError(text string) error {
	return errors.New(prefixes.ServicePrefix + text)
}

var (
	ErrNoEnoughItemsToReserve = serviceError("no enough items to reserve")
	ErrNoEnoughItemsInStock   = serviceError("no enough items in stock")
	ErrAlreadyProcessed       = serviceError("already processed")
)

// После генерации mock-а добавь структуру
// type ExecuteKey struct{}
// и в начало функции WithinTransaction следующее условие, чтобы проходить тесты:
// if ok := arg0.Value(ExecuteKey{}); ok != nil {
//  	return arg1(arg0)
// }

//go:generate mockgen -source=service.go -destination=mocks/service.go
type Interface interface {
	// ChangePriceInStock изменяет цену товара, находящегося в продаже
	ChangePriceInStock(ctx context.Context, data dto.ArticleWithPriceDTO) error
	// Stock возвращает полную информацию о товаре, доступном для продажи, в виде dto.ArticlePriceNameAmount
	Stock(ctx context.Context, data dto.Article) (dto.ArticlePriceNameAmount, error)
	// AddProductToStock добавляет новый товар в ассортимент магазина
	AddProductToStock(ctx context.Context, data dto.ArticlePriceNameAmount) error
	// ChangeAmountInStock изменяет доступное для продажи количество товара
	ChangeAmountInStock(ctx context.Context, data dto.ArticleAmount) error
	// AmountInStock возвращает доступное для продажи количество товара
	AmountInStock(ctx context.Context, data dto.Article) (uint, error)
	// MakeReservation производит резервирование товара для покупателя. Резервирование проводится как для бронирования
	// через интернет, так и во время нахождения товара на кассе (в ожидании оплаты локальным покупателем). В таком случае
	// в качестве номера заказа передаётся номер кассы.
	MakeReservation(ctx context.Context, data dto.ReservationDTO) error
	// CancelReservation снимает бронь с товара/ов
	CancelReservation(ctx context.Context, data dto.Number) error
	// MakeSale уменьшает количества доступного для продажи товара и производит запись в статистику продаж
	MakeSale(ctx context.Context, data []dto.ArticlePriceAmount) error
	// FinishOrder помечает заказ, как выполненный. Данные о содержащихся в заказе товарах переносятся в статистику продаж
	FinishOrder(ctx context.Context, data dto.Number) error
	// TotalSold возвращает количество проданного товара с переданным артикулом за весь период
	TotalSold(ctx context.Context, data dto.Article) (uint, error)
	// TotalSoldInPeriod возвращает количество проданного товара с переданным артикулом за указанный период
	TotalSoldInPeriod(ctx context.Context, data dto.ArticlePeriod) (uint, error)
}
