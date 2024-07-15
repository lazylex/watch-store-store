package dto

import "github.com/lazylex/watch-store-store/internal/dto/validators"

type Price struct {
	Price float64 `json:"price"`
}

// Validate валидация корректности сохраненных в DTO данных.
func (p *Price) Validate() error {
	return validators.Price(p.Price)
}
