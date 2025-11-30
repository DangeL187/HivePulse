package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MessagesReceived = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_received_total",
			Help: "Total messages received by consumer",
		},
	)
	AuthSuccess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_auth_success_total",
			Help: "Messages successfully authenticated",
		},
	)
	AuthFail = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_auth_fail_total",
			Help: "Messages failed authentication",
		},
	)
	MessagesDropped = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_dropped_total",
			Help: "Messages dropped due to full channel",
		},
	)
	ConsumerLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "consumer_latency_seconds",
			Help:    "Latency of message processing in consumer",
			Buckets: prometheus.DefBuckets,
		},
	)
	MessagesSent = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_sent_total",
			Help: "Messages successfully sent to Kafka",
		},
	)
	MessagesSendErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_send_errors_total",
			Help: "Errors while sending messages to Kafka",
		},
	)
)

func RegisterAll() {
	prometheus.MustRegister(MessagesReceived, AuthSuccess, AuthFail,
		MessagesDropped, ConsumerLatency, MessagesSent, MessagesSendErrors)
}
