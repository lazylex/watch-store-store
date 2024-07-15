package dto

import (
	"github.com/lazylex/watch-store-store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store-store/internal/dto/validators"
	"time"
)

type ArticlePriceAmountDate struct {
	Article article.Article `json:"article"`
	Price   float64         `json:"price"`
	Amount  uint            `json:"amount"`
	Date    time.Time       `json:"date"`
}

// Validate валидация корректности сохраненных в DTO данных.
func (h *ArticlePriceAmountDate) Validate() error {
	if err := validators.Article(h.Article); err != nil {
		return err
	}
	if err := validators.Price(h.Price); err != nil {
		return err
	}
	return nil
}
