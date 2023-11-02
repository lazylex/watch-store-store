package request

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/lazylex/watch-store/store/internal/adapters/rest/response"
	"github.com/lazylex/watch-store/store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store/store/internal/dto"
	"github.com/lazylex/watch-store/store/internal/helpers/constants/prefixes"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

const (
	Article  = "article"
	Products = "products[]"
	Amount   = "amount"
	From     = "from"
	To       = "to"
)

// requestErr добавляет к тексту ошибки префикс, указывающий на её принадлежность к запросу
func requestErr(text string) error {
	return errors.New(prefixes.RequestErrorsPrefix + text)
}

var ErrIncorrectDate = requestErr("invalid date passed")
var ErrEmptyProductsData = requestErr("product data not sent")
var ErrEmptyFromDate = requestErr("no 'from' date in request")

// GetArticleUsingChi возвращает артукул продукта
func GetArticleUsingChi(r *http.Request) article.Article {
	return article.Article(chi.URLParam(r, Article))
}

// GetProductDTOs возращает []dto.ProductDTO, полученный из аргументов запроса, содержащего данные о продукте в таком
// виде: products[]=ca-f91w,2100,20&products[]=ca-aw-591,15000,36. Если данных не обнаружено, в заголовок ответа
// записывается http.StatusBadRequest, возвращается пустой слайс и ошибка ErrEmptyProductsData, которая так же
// записывается в лог
func GetProductDTOs(w http.ResponseWriter, r *http.Request, logger *slog.Logger) ([]dto.ProductDTO, error) {
	var products []dto.ProductDTO
	articlePosition, pricePosition, amountPosition := 0, 1, 2
	for _, productsData := range r.URL.Query()[Products] {
		productData := strings.Split(productsData, ",")
		if len(productData) != 3 {
			continue
		}

		price, err := strconv.ParseFloat(productData[pricePosition], 64)
		if err != nil {
			continue
		}

		amount, err := strconv.ParseUint(productData[amountPosition], 10, 64)
		if err != nil {
			continue
		}

		art := article.Article(productData[articlePosition])

		productDTO := dto.ProductDTO{Article: art, Price: price, Amount: uint(amount)}
		if productDTO.Validate() == nil {
			products = append(products, productDTO)
		}
	}

	if len(products) == 0 {
		response.WriteHeaderAndLogAboutErr(w, logger, ErrEmptyProductsData)
		return products, ErrEmptyProductsData
	}

	return products, nil
}
