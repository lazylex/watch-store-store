package validators

import (
	"errors"
	"github.com/lazylex/watch-store-store/internal/domain/aggregates/reservation"
	"github.com/lazylex/watch-store-store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store-store/internal/helpers/constants/prefixes"
	"time"
)

// dtoErr добавляет к тексту ошибки префикс, указывающий на её принадлежность к DTO.
func dtoErr(text string) error {
	return errors.New(prefixes.DTOErrorsPrefix + text)
}

var (
	ErrIncorrectArticle               = dtoErr("incorrect article")
	ErrIncorrectState                 = dtoErr("incorrect state")
	ErrIncorrectOrder                 = dtoErr("incorrect order")
	ErrNegativePrice                  = dtoErr("negative product price")
	ErrZeroPrice                      = dtoErr("zero product price")
	ErrEmptyName                      = dtoErr("empty product name")
	ErrIncorrectDatesOrder            = dtoErr("incorrect dates order")
	ErrDatesIsEqual                   = dtoErr("dates is equal")
	ErrCashRegisterOrder              = dtoErr("order number for cash register can't be higher than 10")
	ErrOrderForInternetCustomer       = dtoErr("number for internet based orders must be higher than 10")
	ErrDuplicateProductsInReservation = dtoErr("duplicate products in reservation")
	ErrNoProductsInReservation        = dtoErr("no products in reservation")
)

// Article функция валидации артикула.
func Article(a article.Article) error {
	r := []rune(a)
	ln := len(r)

	if a == "" || ln > 50 {
		return ErrIncorrectArticle
	}

	if ln > 5 && r[ln-5] == rune('.') {
		switch {
		case r[ln-4]-48 > article.CaseWithHeavyScratches || r[ln-4] < 48:
			return ErrIncorrectArticle
		case r[ln-3]-48 > article.DisplayWithHeavyScratches || r[ln-3] < 48:
			return ErrIncorrectArticle
		case r[ln-2]-48 > article.PackageOpened || r[ln-2] < 48:
			return ErrIncorrectArticle
		case r[ln-1]-48 > article.PackagingWithDamage || r[ln-1] < 48:
			return ErrIncorrectArticle
		}
	}
	return nil
}

// DateOrder функция валидации порядка дат в периоде времени.
func DateOrder(from, to time.Time) error {
	if from.After(to) {
		return ErrIncorrectDatesOrder
	}
	if from.Equal(to) {
		return ErrDatesIsEqual
	}
	return nil
}

// Price функция валидации цены.
func Price(price float64) error {
	if price < 0 {
		return ErrNegativePrice
	}
	if price == 0 {
		return ErrZeroPrice
	}

	return nil
}

// Name функция валидации названия.
func Name(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	return nil
}

// OrderNumber функция валидации номера заказа.
func OrderNumber(order reservation.OrderNumber) error {
	if int(order) <= 0 {
		return ErrIncorrectOrder
	}
	return nil
}
