package reservation

type OrderNumber int

const MaxCashRegisterNumber = 10

const (
	NewForCashRegister = iota + 1
	NewForLocalCustomer
	NewForInternetCustomer
	Finished
	Cancel
)
