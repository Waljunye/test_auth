package metrics

import "github.com/prometheus/client_golang/prometheus"

var HttpRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_request_total",
	Help: "Total number of HTTP requests",
},
	[]string{"method", "status"})
