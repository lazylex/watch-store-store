package handlers

import (
	"net/http"
)

type Interface interface {
	StockRecord(w http.ResponseWriter, r *http.Request)
	AmountInStock(w http.ResponseWriter, r *http.Request)
	UpdatePriceInStock(w http.ResponseWriter, r *http.Request)
	UpdateAmountInStock(w http.ResponseWriter, r *http.Request)
	AddToStock(w http.ResponseWriter, r *http.Request)
	SoldAmount(w http.ResponseWriter, r *http.Request)
	MakeReservation(w http.ResponseWriter, r *http.Request)
	CancelReservation(w http.ResponseWriter, r *http.Request)
	MakeLocalSale(w http.ResponseWriter, r *http.Request)
	FinishOrder(w http.ResponseWriter, r *http.Request)
}
