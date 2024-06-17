package dto

import (
	"errors"
	"github.com/lazylex/watch-store-store/internal/dto/validators"
	"testing"
)

func TestOrderNumberDTO(t *testing.T) {
	t.Run("negative order number", func(t *testing.T) {
		on := OrderNumberDTO{OrderNumber: -1}
		err := on.Validate()
		if !errors.Is(err, validators.ErrIncorrectOrder) {
			t.Fail()
		}
	})
}
