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

	sshDisconnect = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ssh_disconnect_total",
			Help: "SSH disconnect events",
		},
		[]string{"user"},
	)

	sshLoginLast = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ssh_login_last_timestamp",
			Help: "Last SSH login timestamp (unix)",
		},
		[]string{"user", "ip"},
	)

	sshSuLast = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ssh_su_last_timestamp",
			Help: "Last su event timestamp (unix)",
		},
		[]string{"from_user", "to_user"},
	)

	sshSessionLast = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ssh_session_last_timestamp",
			Help: "Last session event timestamp (unix)",
		},
		[]string{"type", "user"},
	)

	sshDisconnectLast = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ssh_disconnect_last_timestamp",
			Help: "Last disconnect timestamp (unix)",
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
	prometheus.MustRegister(sshLoginLast)
	prometheus.MustRegister(sshSuLast)
	prometheus.MustRegister(sshSessionLast)
	prometheus.MustRegister(sshDisconnectLast)
}
