package dto

import (
	"errors"
	"github.com/lazylex/watch-store-store/internal/dto/validators"
	"testing"
)

func TestPriceDTO(t *testing.T) {
	t.Run("zero price", func(t *testing.T) {
		p := Price{Price: 0}
		err := p.Validate()
		if !errors.Is(err, validators.ErrZeroPrice) {
			t.Fail()
		}
	})

	t.Run("negative price", func(t *testing.T) {
		p := Price{Price: -5}
		err := p.Validate()
		if !errors.Is(err, validators.ErrNegativePrice) {
			t.Fail()
		}
	})
}
