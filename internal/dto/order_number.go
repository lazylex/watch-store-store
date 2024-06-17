package dto

import (
	"github.com/lazylex/watch-store-store/internal/domain/aggregates/reservation"
	"github.com/lazylex/watch-store-store/internal/dto/validators"
)

type OrderNumberDTO struct {
	OrderNumber reservation.OrderNumber `json:"order_number"`
}

// Validate валидация корректности сохраненных в DTO данных.
func (on *OrderNumberDTO) Validate() error {
	return validators.OrderNumber(on.OrderNumber)
}
