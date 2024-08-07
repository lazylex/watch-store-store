package dto

import (
	"github.com/lazylex/watch-store-store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store-store/internal/dto/validators"
)

type ArticlePrice struct {
	Article article.Article `json:"article"`
	Price   float64         `json:"price"`
}

// Validate валидация корректности сохраненных в DTO данных.
func (a *ArticlePrice) Validate() error {
	if err := validators.Price(a.Price); err != nil {
		return err
	}
	if err := validators.Article(a.Article); err != nil {
		return err
	}
	return nil
}
