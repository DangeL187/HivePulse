package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MessagesConsumed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "consumer_messages_total",
			Help: "Total number of messages consumed",
		},
	)
	ConsumerLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "consumer_latency_seconds",
			Help:    "Latency of message processing in consumer",
			Buckets: prometheus.DefBuckets,
		},
	)

	BatchesFlushed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "flusher_batches_total",
			Help: "Total number of message batches flushed",
		},
	)
	MessagesFlushed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "flusher_messages_total",
			Help: "Total number of messages flushed",
		},
	)
	FlushDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "flusher_duration_seconds",
			Help:    "Duration of flush operations",
			Buckets: prometheus.DefBuckets,
		},
	)
	FlushErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "flusher_errors_total",
			Help: "Total number of errors in flusher",
		},
	)
)

func RegisterAll() {
	prometheus.MustRegister(MessagesConsumed, ConsumerLatency, BatchesFlushed, MessagesFlushed, FlushDuration,
		FlushErrors)
}
