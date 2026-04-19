package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Login Events
	sshLogins = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ssh_logins_total",
			Help: "SSH login attempts (success/fail)",
		},
		[]string{"status", "user", "ip"},
	)

	// Session Events (nur user)
	sshSessionOpen = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ssh_session_open_total",
			Help: "SSH session opened",
		},
		[]string{"user"},
	)

	sshSessionClose = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ssh_session_close_total",
			Help: "SSH session closed",
		},
		[]string{"user"},
	)

	// SU Events (nur hier from/to!)
	sshSuOpen = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ssh_su_open_total",
			Help: "su session opened",
		},
		[]string{"from_user", "to_user"},
	)

	sshSuClose = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ssh_su_close_total",
			Help: "su session closed",
		},
		[]string{"user"},
	)

	// Disconnect / sonstige Events
	sshDisconnect = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ssh_disconnect_total",
			Help: "SSH disconnect events",
		},
		[]string{"user"},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(sshLogins)
	prometheus.MustRegister(sshSessionOpen)
	prometheus.MustRegister(sshSessionClose)
	prometheus.MustRegister(sshSuOpen)
	prometheus.MustRegister(sshSuClose)
	prometheus.MustRegister(sshDisconnect)
}
