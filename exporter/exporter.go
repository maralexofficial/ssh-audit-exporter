package exporter

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Rule struct {
	Name   string
	Regex  string
	Type   string
	Metric string
}

type Config struct {
	Rules []Rule
}

var defaultConfig = Config{
	Rules: []Rule{
		{
			Name:   "ssh_login_success",
			Type:   "success",
			Metric: "ssh_logins",
			Regex:  `Accepted .* for (\w+) from ([0-9.]+)`,
		},

		{
			Name:   "ssh_login_fail",
			Type:   "fail",
			Metric: "ssh_logins",
			Regex:  `Failed password for (\w+) from ([0-9.]+)`,
		},

		{
			Name:   "ssh_auth_failure",
			Type:   "fail",
			Metric: "ssh_events",
			Regex:  `authentication failure;.*user=(\w+) from ([0-9.]+)`,
		},

		{
			Name:   "ssh_connection_closed",
			Type:   "info",
			Metric: "ssh_events",
			Regex:  `Connection closed by .* user (\w+) from ([0-9.]+)`,
		},

		{
			Name:   "ssh_session_open",
			Type:   "info",
			Metric: "ssh_sessions",
			Regex:  `session opened for user (\w+) from ([0-9.]+)`,
		},

		{
			Name:   "ssh_session_close",
			Type:   "info",
			Metric: "ssh_sessions",
			Regex:  `session closed for user (\w+) from ([0-9.]+)`,
		},
	},
}

func ParseRules(input []string) ([]Rule, error) {

	if len(input) == 0 {
		return defaultConfig.Rules, nil
	}

	var rules []Rule

	for _, r := range input {

		parts := strings.SplitN(r, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid rule format: %s (expected name:regex)", r)
		}

		name := strings.TrimSpace(parts[0])
		regex := strings.TrimSpace(parts[1])

		rule := Rule{
			Name:   name,
			Regex:  regex,
			Type:   "any",
			Metric: "ssh_events",
		}

		switch {
		case strings.Contains(name, "success"):
			rule.Type = "success"
			rule.Metric = "ssh_logins"

		case strings.Contains(name, "fail"):
			rule.Type = "fail"
			rule.Metric = "ssh_logins"

		case strings.Contains(name, "su"):
			rule.Type = "info"
			rule.Metric = "ssh_sessions"
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

func LoadConfig(cliSuccess, cliFail string) Config {
	cfg := defaultConfig

	if env := os.Getenv("SSH_SUCCESS_REGEX"); env != "" {
		cfg.Rules[0].Regex = env
	}

	if env := os.Getenv("SSH_FAIL_REGEX"); env != "" {
		cfg.Rules[1].Regex = env
	}

	if cliSuccess != "" {
		cfg.Rules[0].Regex = cliSuccess
	}

	if cliFail != "" {
		cfg.Rules[1].Regex = cliFail
	}

	return cfg
}

type compiledRule struct {
	Rule  Rule
	Regex *regexp.Regexp
}

type Parser struct {
	rules []compiledRule
}

func NewParser(rules []Rule) *Parser {

	compiled := make([]compiledRule, 0, len(rules))

	for _, r := range rules {
		re, err := regexp.Compile(r.Regex)
		if err != nil {
			panic(err)
		}

		compiled = append(compiled, compiledRule{
			Rule:  r,
			Regex: re,
		})
	}

	return &Parser{
		rules: compiled,
	}
}

func (p *Parser) Parse(line string) {

	for _, r := range p.rules {

		m := r.Regex.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		switch r.Rule.Name {

		case "ssh_success":
			user := m[1]
			ip := m[2]
			sshLogins.WithLabelValues("success", user, ip).Inc()

		case "ssh_failed":
			user := m[1]
			ip := m[2]
			sshLogins.WithLabelValues("failed", user, ip).Inc()

		case "su_session":
			action := m[1]
			user := m[2]
			sshSessions.WithLabelValues(action, user).Inc()

		case "disconnect":
			user := m[1]
			sshEvents.WithLabelValues("disconnect", user).Inc()

		default:
			sshEvents.WithLabelValues(r.Rule.Name, "unknown").Inc()
		}

		return
	}
}
