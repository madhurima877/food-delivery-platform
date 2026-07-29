package metrics

import "github.com/prometheus/client_golang/prometheus"

var OrdersCreated = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "orders_created_total",
		Help: "Total number of orders created",
	},
)

func init() {
	prometheus.MustRegister(OrdersCreated)
	prometheus.MustRegister(OrderCreationErrors)
}

var OrderCreationErrors = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "order_creation_errors_total",
		Help: "Total number of order creation errors",
	},
)
