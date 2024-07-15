package dto

import (
	"errors"
	"github.com/lazylex/watch-store-store/internal/dto/validators"
	"testing"
)

func TestProduct(t *testing.T) {
	t.Run("error price", func(t *testing.T) {
		p := ArticlePriceAmount{Article: "test-9", Price: 0, Amount: 10}
		err := p.Validate()
		if !errors.Is(err, validators.ErrZeroPrice) {
			t.Fail()
		}
	})
}
