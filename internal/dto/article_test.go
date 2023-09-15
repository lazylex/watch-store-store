package dto

import (
	"errors"
	"github.com/lazylex/watch-store/store/internal/dto/validators"
	"testing"
)

func TestArticle(t *testing.T) {
	a := ArticleDTO{Article: "test-9.0009"}
	t.Run("incorrect article", func(t *testing.T) {
		err := a.Validate()
		if !errors.Is(err, validators.ErrIncorrectArticle) {
			t.Fail()
		}
	})
}
