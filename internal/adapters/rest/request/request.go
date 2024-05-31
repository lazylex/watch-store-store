package request

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/lazylex/watch-store/store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store/store/internal/helpers/constants/prefixes"
	"net/http"
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

// ArticleUsingChi возвращает артикул продукта.
func ArticleUsingChi(r *http.Request) article.Article {
	return article.Article(chi.URLParam(r, Article))
}
