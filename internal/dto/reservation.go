package dto

import (
	rs "github.com/lazylex/watch-store/store/internal/domain/aggregates/reservation"
	"github.com/lazylex/watch-store/store/internal/dto/validators"
	"time"
)

type ReservationDTO struct {
	Products    []ProductDTO
	OrderNumber rs.OrderNumber `json:"order_number"`
	Date        time.Time      `json:"date"`
	State       uint           `json:"state"`
}

// Validate валидация корректности сохраненных в DTO данных
func (r *ReservationDTO) Validate() error {
	if err := validators.OrderNumber(r.OrderNumber); err != nil {
		return err
	}

	if r.State != rs.NewForLocalCustomer && r.State != rs.NewForInternetCustomer && r.State != rs.NewForCashRegister {
		return validators.ErrIncorrectState
	}

	if r.State == rs.NewForCashRegister && r.OrderNumber > 10 {
		return validators.ErrCashRegisterOrder
	}

	if r.State != rs.NewForCashRegister && r.OrderNumber <= 10 {
		return validators.ErrOrderForInternetCustomer
	}

	for _, product := range r.Products {
		if err := product.Validate(); err != nil {
			return err
		}
	}

	return nil
}
