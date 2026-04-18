package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var SSHLogins = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "ssh_logins_total",
		Help: "Total number of SSH login attempts",
	},
	[]string{"status"},
)

func RegisterMetrics() {
	prometheus.MustRegister(SSHLogins)
}