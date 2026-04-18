package exporter

import (
	"os"
	"regexp"
)

type Config struct {
	SuccessRegex string
	FailRegex    string
}

var defaultConfig = Config{
	SuccessRegex: `Accepted .* for (\w+) from ([0-9.]+)`,
	FailRegex:    `Failed .* for (\w+) from ([0-9.]+)`,
}

func LoadConfig(cliSuccess, cliFail string) Config {
	cfg := defaultConfig

	if env := os.Getenv("SSH_SUCCESS_REGEX"); env != "" {
		cfg.SuccessRegex = env
	}

	if env := os.Getenv("SSH_FAIL_REGEX"); env != "" {
		cfg.FailRegex = env
	}

	if cliSuccess != "" {
		cfg.SuccessRegex = cliSuccess
	}

	if cliFail != "" {
		cfg.FailRegex = cliFail
	}

	return cfg
}

func CompileRegex(cfg Config) (*regexp.Regexp, *regexp.Regexp, error) {

	success, err := regexp.Compile(cfg.SuccessRegex)
	if err != nil {
		return nil, nil, err
	}

	fail, err := regexp.Compile(cfg.FailRegex)
	if err != nil {
		return nil, nil, err
	}

	return success, fail, nil
}

var (
	successRegex *regexp.Regexp
	failRegex    *regexp.Regexp
)

func InitParser(success, fail *regexp.Regexp) {
	successRegex = success
	failRegex = fail
}

func ParseLine(line string) {

	if successRegex != nil {
		if matches := successRegex.FindStringSubmatch(line); matches != nil {
			user := matches[1]
			ip := matches[2]

			sshLogins.WithLabelValues("success", user, ip).Inc()
			return
		}
	}

	if failRegex != nil {
		if matches := failRegex.FindStringSubmatch(line); matches != nil {
			user := matches[1]
			ip := matches[2]

			sshLogins.WithLabelValues("failed", user, ip).Inc()
			return
		}
	}
}