package service

//go:generate mockgen -source=service.go -destination=mocks/service.go
type MetricsInterface interface {
	CancelOrdersInc()
	PlacedInternetOrdersInc()
	PlacedLocalOrdersInc()
}
