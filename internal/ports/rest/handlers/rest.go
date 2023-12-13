package handlers

import (
	"net/http"
)

type Interface interface {
	GetStockRecord(w http.ResponseWriter, r *http.Request)
	GetAmountInStock(w http.ResponseWriter, r *http.Request)
	UpdatePriceInStock(w http.ResponseWriter, r *http.Request)
	UpdateAmountInStock(w http.ResponseWriter, r *http.Request)
	AddToStock(w http.ResponseWriter, r *http.Request)
	GetSoldAmount(w http.ResponseWriter, r *http.Request)
	MakeReservation(w http.ResponseWriter, r *http.Request)
	CancelReservation(w http.ResponseWriter, r *http.Request)
	MakeLocalSale(w http.ResponseWriter, r *http.Request)
	FinishOrder(w http.ResponseWriter, r *http.Request)
	GetAllHandlers(w http.ResponseWriter, r *http.Request)
}
