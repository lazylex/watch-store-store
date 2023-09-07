package dto

import (
	"github.com/lazylex/watch-store/store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store/store/internal/dto/validators"
)

type ProductDTO struct {
	Article article.Article `json:"article"`
	Price   float64         `json:"price"`
	Amount  uint            `json:"amount"`
}

// Validate валидация корректности сохраненных в DTO данных
func (p *ProductDTO) Validate() error {
	if err := validators.Article(p.Article); err != nil {
		return err
	}
	if err := validators.Price(p.Price); err != nil {
		return err
	}
	return nil
}
