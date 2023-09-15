package dto

import (
	"errors"
	"github.com/lazylex/watch-store/store/internal/dto/validators"
	"testing"
)

func TestNamedProductDTO(t *testing.T) {
	t.Run("incorrect article", func(t *testing.T) {
		a := NamedProductDTO{Article: "test-9.---9", Price: 1000, Name: "ok name"}
		if !errors.Is(a.Validate(), validators.ErrIncorrectArticle) {
			t.Fail()
		}
	})

	t.Run("zero price", func(t *testing.T) {
		a := NamedProductDTO{Article: "test-9", Price: 0, Name: "ok name"}
		if !errors.Is(a.Validate(), validators.ErrZeroPrice) {
			t.Fail()
		}
	})

	t.Run("negative price", func(t *testing.T) {
		a := NamedProductDTO{Article: "test-9", Price: -10, Name: "ok name"}
		if !errors.Is(a.Validate(), validators.ErrNegativePrice) {
			t.Fail()
		}
	})

	t.Run("empty name", func(t *testing.T) {
		a := NamedProductDTO{Article: "test-9", Price: 1000, Name: ""}
		if !errors.Is(a.Validate(), validators.ErrEmptyName) {
			t.Fail()
		}
	})

	t.Run("correct article with price", func(t *testing.T) {
		a := NamedProductDTO{Article: "test-9", Price: 1000, Name: "ok name"}
		if a.Validate() != nil {
			t.Fail()
		}
	})
}
