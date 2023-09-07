package product

import (
	"github.com/lazylex/watch-store/store/internal/domain/value_objects/article"
)

type Product struct {
	article article.Article
	price   float64
	amount  uint
}
