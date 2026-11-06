package booking

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Tracks total HTTP requests
	HttpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cinema_http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	// Tracks seats booked per movie
	SeatsBookedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cinema_seats_booked_total",
		Help: "Total number of seats successfully booked",
	}, []string{"movie_id"})
)