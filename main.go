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

var cliLogFile string
var cliSuccessRegex string
var cliFailRegex string

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
	successRegex, failRegex, _ := exporter.CompileRegex(cfg)
	exporter.InitParser(successRegex, failRegex)

	mode := source.GetSourceType()

	switch mode {

	case source.Journal:
		logger.Info("Using journald as source")
		go journal.TailSSHJournal(exporter.ParseLine)

	default:
		logFile, err := logfile.GetLogFile(cliLogFile)
		if err != nil {
			logger.Error("No log file found")
			return
		}

		logger.Success("Using file: " + logFile)
		go logfile.TailFile(logFile, exporter.ParseLine)
	}

	http.Handle("/metrics", promhttp.Handler())

	logger.Info("Metrics available on :9100")
	http.ListenAndServe(":9100", nil)
}