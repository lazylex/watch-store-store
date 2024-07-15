package dto

import (
	"github.com/lazylex/watch-store-store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store-store/internal/dto/validators"
)

type ArticlePriceNameAmount struct {
	Name    string          `json:"name"`
	Article article.Article `json:"article"`
	Price   float64         `json:"price"`
	Amount  uint            `json:"amount"`
}

// Validate валидация корректности сохраненных в DTO данных.
func (np *ArticlePriceNameAmount) Validate() error {
	if err := validators.Name(np.Name); err != nil {
		return err
	}
	if err := validators.Article(np.Article); err != nil {
		return err
	}
	if err := validators.Price(np.Price); err != nil {
		return err
	}

	return nil
}
