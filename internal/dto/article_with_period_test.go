package dto

import (
	"errors"
	"github.com/lazylex/watch-store/store/internal/dto/validators"
	"testing"
	"time"
)

func TestArticleWithPeriodDTO(t *testing.T) {
	t.Run("incorrect article", func(t *testing.T) {
		a := ArticleWithPeriodDTO{Article: "test-9.....", From: time.Now(), To: time.Now()}
		err := a.Validate()
		if !errors.Is(err, validators.ErrIncorrectArticle) {
			t.Fail()
		}
	})

	t.Run("incorrect dates order", func(t *testing.T) {
		a := ArticleWithPeriodDTO{Article: "test-9",
			From: time.Date(2021, 1, 1, 1, 1, 1, 0, time.Local),
			To:   time.Date(2020, 1, 1, 1, 1, 1, 0, time.Local)}
		err := a.Validate()
		if !errors.Is(err, validators.ErrIncorrectDatesOrder) {
			t.Fail()
		}
	})

	t.Run("correct", func(t *testing.T) {
		a := ArticleWithPeriodDTO{Article: "test-9",
			From: time.Date(2020, 1, 1, 1, 1, 1, 0, time.Local),
			To:   time.Date(2021, 1, 1, 1, 1, 1, 0, time.Local)}

		if a.Validate() != nil {
			t.Fail()
		}
	})
}
