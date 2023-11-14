package metrics

import "github.com/prometheus/client_golang/prometheus"

type Service struct {
	canceledOrders *prometheus.CounterVec
}

func (s *Service) CancelOrdersInc() {
	s.canceledOrders.With(prometheus.Labels{}).Inc()
}

// createCanceledOrdersTotalMetric создает и регистрирует метрику canceled_orders_total, являющуюся счетчиком отмененных
// запросов
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
