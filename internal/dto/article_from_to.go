package dto

import (
	"github.com/lazylex/watch-store-store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store-store/internal/dto/validators"
	"time"
)

type ArticleFromTo struct {
	Article article.Article `json:"article"`
	From    time.Time       `json:"from"`
	To      time.Time       `json:"to"`
}

// Validate валидация корректности сохраненных в DTO данных.
func (ap *ArticleFromTo) Validate() error {
	if err := validators.Article(ap.Article); err != nil {
		return err
	}
	if err := validators.DateOrder(ap.From, ap.To); err != nil {
		return err
	}
	return nil
}
