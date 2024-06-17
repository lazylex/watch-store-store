package sales_history

import (
	"github.com/lazylex/watch-store-store/internal/domain/value_objects/product"
	"time"
)

type SalesHistory struct {
	product product.Product
	date    time.Time
}
