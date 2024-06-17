package dto

import (
	rs "github.com/lazylex/watch-store-store/internal/domain/aggregates/reservation"
	"github.com/lazylex/watch-store-store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store-store/internal/dto/validators"
	"time"
)

type ReservationDTO struct {
	Products    []ProductDTO   `json:"products"`
	OrderNumber rs.OrderNumber `json:"order_number"`
	Date        time.Time      `json:"date"`
	State       uint           `json:"state"`
}

// IsNew возвращает true, если бронь еще не была снята (по причине отмены или завершения заказа).
func (r *ReservationDTO) IsNew() bool {
	return r.State == rs.NewForInternetCustomer || r.State == rs.NewForLocalCustomer || r.State == rs.NewForCashRegister
}

// Validate валидация корректности сохраненных в DTO данных.
func (r *ReservationDTO) Validate() error {
	if err := validators.OrderNumber(r.OrderNumber); err != nil {
		return err
	}

	if !r.IsNew() {
		return validators.ErrIncorrectState
	}

	if r.State == rs.NewForCashRegister && r.OrderNumber > rs.MaxCashRegisterNumber {
		return validators.ErrCashRegisterOrder
	}

	if r.State != rs.NewForCashRegister && r.OrderNumber <= rs.MaxCashRegisterNumber {
		return validators.ErrOrderForInternetCustomer
	}

	if len(r.Products) == 0 {
		return validators.ErrNoProductsInReservation
	}

	articles := make(map[article.Article]struct{})
	for _, product := range r.Products {
		if err := product.Validate(); err != nil {
			return err
		}
		if _, ok := articles[product.Article]; !ok {
			articles[product.Article] = struct{}{}
		} else {
			return validators.ErrDuplicateProductsInReservation
		}
	}

	return nil
}
