package main

import (
	"fmt"
	"net/http"
	"os"

	"ssh-audit-exporter/exporter"
	"ssh-audit-exporter/logger"

	"github.com/hpcloud/tail"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
)

var logPaths = []string{
	"/var/log/auth.log",
	"/var/log/secure",
}

func findLogFile() (string, error) {
	for _, path := range logPaths {
		file, err := os.Open(path)
		if err == nil {
			file.Close()
			return path, nil
		}
	}
	return "", fmt.Errorf("no readable SSH log file found (missing or permission denied)")
}

func tailLog(file string) {
	logger.Info("Starting log monitoring: " + file)

	t, err := tail.TailFile(file, tail.Config{
		Follow: true,
		ReOpen: true,
	})
	if err != nil {
		logger.Error("Failed to open log file: " + err.Error())
		os.Exit(1)
	}

	for line := range t.Lines {
		if line == nil {
			continue
		}
		exporter.ParseLine(line.Text)
	}
}

func main() {
	logger.Info("Starting SSH audit exporter")

	exporter.RegisterMetrics()

	logFile, err := findLogFile()
	if err != nil {
		logger.Error(err.Error())
		logger.Warning("Checked log paths:")
		for _, p := range logPaths {
			logger.Warning(" - " + p)
		}
		logger.Warning("System might be using journald instead of log files")
		return
	}

	logger.Success("Using log file: " + logFile)

	go tailLog(logFile)

	http.Handle("/metrics", promhttp.Handler())

	logger.Info("Metrics available at http://localhost:9100/metrics")

	err = http.ListenAndServe(":9100", nil)
	if err != nil {
		logger.Error("HTTP server failed: " + err.Error())
	}
}