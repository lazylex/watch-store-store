package request

import (
	"errors"
	"github.com/lazylex/watch-store/store/internal/helpers/constants/prefixes"
)

const (
	Article = "article"
	Amount  = "amount"
	From    = "from"
	To      = "to"
)

// requestErr добавляет к тексту ошибки префикс, указывающий на её принадлежность к запросу.
func requestErr(text string) error {
	return errors.New(prefixes.RequestErrorsPrefix + text)
}

var ErrIncorrectDate = requestErr("invalid date passed")
var ErrEmptyFromDate = requestErr("no 'from' date in request")
