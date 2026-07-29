package metrics

import "github.com/prometheus/client_golang/prometheus"

var InventoryReservations = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "inventory_reservations_total",
		Help: "Total number of successful inventory reservations",
	},
)

var InventoryFailures = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "inventory_failures_total",
		Help: "Total number of failed inventory reservations",
	},
)

func init() {
	prometheus.MustRegister(InventoryReservations)
	prometheus.MustRegister(InventoryFailures)
}
