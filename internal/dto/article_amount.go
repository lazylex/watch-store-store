package dto

import (
	"github.com/lazylex/watch-store-store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store-store/internal/dto/validators"
)

type ArticleAmount struct {
	Article article.Article `json:"article"`
	Amount  uint            `json:"amount"`
}

// Validate валидация корректности сохраненных в DTO данных.
func (a *ArticleAmount) Validate() error {
	// Amount не валидируем, нулевое значение считаем валидным
	return validators.Article(a.Article)
}
