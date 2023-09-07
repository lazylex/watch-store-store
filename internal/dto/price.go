package dto

import "github.com/lazylex/watch-store/store/internal/dto/validators"

type PriceDTO struct {
	Price float64 `json:"price"`
}

// Validate валидация корректности сохраненных в DTO данных
func (p *PriceDTO) Validate() error {
	return validators.Price(p.Price)
}
