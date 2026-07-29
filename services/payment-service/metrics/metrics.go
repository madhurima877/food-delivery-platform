package metrics

import "github.com/prometheus/client_golang/prometheus"

var PaymentsCompleted = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "payments_completed_total",
		Help: "Total number of completed payments",
	},
)
var PaymentsFailed = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "payments_failed_total",
		Help: "Total number of failed payments",
	},
)

func init() {
	prometheus.MustRegister(PaymentsCompleted)
	prometheus.MustRegister(PaymentsFailed)
	prometheus.MustRegister(PaymentProcessingDuration)
}

var PaymentProcessingDuration = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name: "payment_processing_duration_seconds",
		Help: "Time taken to process payments",
	},
)
