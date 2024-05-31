package metrics

import "github.com/prometheus/client_golang/prometheus"

type Service struct {
	canceledOrders       *prometheus.CounterVec
	placedInternetOrders *prometheus.CounterVec
	placedLocalOrders    *prometheus.CounterVec
}

// CancelOrdersInc увеличивает счетчик отмененных заказов.
func (s *Service) CancelOrdersInc() {
	s.canceledOrders.With(prometheus.Labels{}).Inc()
}

// PlacedInternetOrdersInc увеличивает счетчик размещенных заказов интернет-магазина с доставкой к покупателю.
func (s *Service) PlacedInternetOrdersInc() {
	s.placedInternetOrders.With(prometheus.Labels{}).Inc()
}

// PlacedLocalOrdersInc увеличивает счетчик размещенных заказов для товаров, бронированных для самовывоза из магазина.
func (s *Service) PlacedLocalOrdersInc() {
	s.placedLocalOrders.With(prometheus.Labels{}).Inc()
}

// createCanceledOrdersTotalMetric создает и регистрирует метрику canceled_orders_total, являющуюся счетчиком отмененных
// заказов.
func createCanceledOrdersTotalMetric() (*prometheus.CounterVec, error) {
	var err error
	orders := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "canceled_orders_total",
		Namespace: NAMESPACE,
		Help:      "Count of canceled orders",
	}, []string{})
	if err = prometheus.Register(orders); err != nil {
		return nil, err
	}

	orders.With(prometheus.Labels{})

	return orders, nil
}

// createPlacedInternetOrdersTotalMetric создает и регистрирует метрику placed_internet_orders_total, являющуюся
// счетчиком совершенных заказов интернет-магазина для доставки покупателю.
func createPlacedInternetOrdersTotalMetric() (*prometheus.CounterVec, error) {
	var err error
	orders := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "placed_internet_orders_total",
		Namespace: NAMESPACE,
		Help:      "Count of internet placed orders",
	}, []string{})
	if err = prometheus.Register(orders); err != nil {
		return nil, err
	}

	orders.With(prometheus.Labels{})

	return orders, nil
}

// createPlacedLocalOrdersTotalMetric создает и регистрирует метрику placed_local_orders_total, являющуюся
// счетчиком совершенных заказов для самовывоза покупателем.
func createPlacedLocalOrdersTotalMetric() (*prometheus.CounterVec, error) {
	var err error
	orders := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "placed_local_orders_total",
		Namespace: NAMESPACE,
		Help:      "Count of local placed orders",
	}, []string{})
	if err = prometheus.Register(orders); err != nil {
		return nil, err
	}

	orders.With(prometheus.Labels{})

	return orders, nil
}
