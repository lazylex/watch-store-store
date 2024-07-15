package dto

import (
	"github.com/lazylex/watch-store-store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store-store/internal/dto/validators"
)

type Article struct {
	Article article.Article `json:"article"`
}

// Validate валидация корректности сохраненных в DTO данных.
func (a *Article) Validate() error {
	return validators.Article(a.Article)
}
