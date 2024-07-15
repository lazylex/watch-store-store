package dto

import (
	"github.com/lazylex/watch-store-store/internal/domain/aggregates/reservation"
	"github.com/lazylex/watch-store-store/internal/dto/validators"
)

type Number struct {
	OrderNumber reservation.OrderNumber `json:"order_number"`
}

// Validate валидация корректности сохраненных в DTO данных.
func (on *Number) Validate() error {
	return validators.OrderNumber(on.OrderNumber)
}
