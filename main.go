package main

import (
	"flag"
	"net/http"
	"strings"

	"ssh-audit-exporter/exporter"
	"ssh-audit-exporter/internal/journal"
	"ssh-audit-exporter/internal/logfile"
	"ssh-audit-exporter/internal/source"
	"ssh-audit-exporter/logger"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cliLogFile string
	cliRules   ruleFlags
)

type ruleFlags []string

func (r *ruleFlags) String() string {
	return strings.Join(*r, ",")
}

func (r *ruleFlags) Set(value string) error {
	*r = append(*r, value)
	return nil
}

func init() {
	flag.StringVar(&cliLogFile, "logfile", "", "path to log file")
	flag.Var(&cliRules, "rule", "filter rules in format type:regex (can be repeated)")
}

func main() {

	_ = godotenv.Load()
	flag.Parse()

	logger.Info("Starting SSH audit exporter")

	exporter.RegisterMetrics()

	rules, err := exporter.ParseRules(cliRules)
	if err != nil {
		logger.Error("invalid rules: " + err.Error())
		return
	}

	parser := exporter.NewParser(rules)

	mode := source.GetSourceType()

	switch mode {

	case source.Journal:
		logger.Info("Using journald (journalctl stream) as source")

		go func() {
			if err := journal.TailSSHJournal(parser.Parse); err != nil {
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
			if err := logfile.TailFile(logFile, parser.Parse); err != nil {
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
