package metrics

import "github.com/prometheus/client_golang/prometheus"

func New() {
	prometheus.MustRegister(HttpRequestTotal)
}
