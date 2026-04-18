package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var sshLogins = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "ssh_logins_total",
		Help: "Total number of SSH login attempts",
	},
	[]string{"status", "user", "ip"},
)

func RegisterMetrics() {
	prometheus.MustRegister(sshLogins)
}