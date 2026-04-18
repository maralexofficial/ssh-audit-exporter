package main

import (
	"flag"
	"net/http"

	"ssh-audit-exporter/exporter"
	"ssh-audit-exporter/internal/journal"
	"ssh-audit-exporter/internal/logfile"
	"ssh-audit-exporter/internal/source"
	"ssh-audit-exporter/logger"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cliLogFile      string
	cliSuccessRegex string
	cliFailRegex    string
)

func init() {
	flag.StringVar(&cliLogFile, "logfile", "", "")
	flag.StringVar(&cliSuccessRegex, "success-regex", "", "")
	flag.StringVar(&cliFailRegex, "fail-regex", "", "")
}

func main() {

	_ = godotenv.Load()
	flag.Parse()

	logger.Info("Starting SSH audit exporter")

	exporter.RegisterMetrics()

	cfg := exporter.LoadConfig(cliSuccessRegex, cliFailRegex)

	successRegex, failRegex, err := exporter.CompileRegex(cfg)
	if err != nil {
		logger.Error("Invalid regex configuration: " + err.Error())
		return
	}

	exporter.InitParser(successRegex, failRegex)

	mode := source.GetSourceType()

	switch mode {

	case source.Journal:
		logger.Info("Using journald (journalctl stream) as source")

		go func() {
			if err := journal.TailSSHJournal(exporter.ParseLine); err != nil {
				logger.Error("journalctl stream failed: " + err.Error())
			}
		}()

	default:
		logFile, err := logfile.GetLogFile(cliLogFile)
		if err != nil {
			logger.Error("No log file found: " + err.Error())
			return
		}

		logger.Success("Using file: " + logFile)

		go func() {
			if err := logfile.TailFile(logFile, exporter.ParseLine); err != nil {
				logger.Error("file tail failed: " + err.Error())
			}
		}()
	}

	http.Handle("/metrics", promhttp.Handler())

	logger.Info("Metrics available on :9100")

	if err := http.ListenAndServe(":9100", nil); err != nil {
		logger.Error("HTTP server failed: " + err.Error())
	}
}