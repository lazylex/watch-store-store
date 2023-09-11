package reservation

import (
	"github.com/lazylex/watch-store/store/internal/domain/value_objects/product"
	"time"
)

type OrderNumber int

const MaxCashRegisterNumber = 10

const (
	NewForCashRegister  = iota
	NewForLocalCustomer = iota
	NewForInternetCustomer
	Finished
	Cancel
)

type Reservation struct {
	product     []*product.Product
	date        time.Time
	orderNumber OrderNumber
	state       uint
}
