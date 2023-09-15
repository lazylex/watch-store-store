package dto

import (
	"errors"
	"github.com/lazylex/watch-store/store/internal/dto/validators"
	"testing"
)

func TestArticleWithAmountDTO(t *testing.T) {
	a := ArticleWithAmountDTO{Article: "test-9.0009", Amount: 0}
	t.Run("incorrect article", func(t *testing.T) {
		err := a.Validate()
		if !errors.Is(err, validators.ErrIncorrectArticle) {
			t.Fail()
		}
	})
}
