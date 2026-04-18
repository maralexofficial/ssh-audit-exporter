package exporter

import "regexp"

var (
	successRegex = regexp.MustCompile(`Accepted .* for (\w+)`)
	failRegex    = regexp.MustCompile(`Failed .* for`)
)

func ParseLine(line string) {
	if successRegex.MatchString(line) {
		SSHLogins.WithLabelValues("success").Inc()
		return
	}

	if failRegex.MatchString(line) {
		SSHLogins.WithLabelValues("failed").Inc()
		return
	}
}