package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	sshLogins = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ssh_logins_total",
			Help: "SSH login attempts (success/fail)",
		},
		[]string{"status", "user", "ip"},
	)

	sshSessions = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ssh_sessions_total",
			Help: "SSH session events (su, session open/close)",
		},
		[]string{"action", "user"},
	)

	sshEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ssh_events_total",
			Help: "Generic SSH-related events",
		},
		[]string{"type", "user"},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(sshLogins)
	prometheus.MustRegister(sshSessions)
	prometheus.MustRegister(sshEvents)
}
