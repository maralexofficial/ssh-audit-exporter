package exporter

import (
	"fmt"
	"regexp"
	"strings"
)

// -----------------------------
// Rule Definition
// -----------------------------

type Rule struct {
	Name   string
	Regex  string
	Type   string // success | fail | info | any
	Metric string
	Labels []string
}

type compiledRule struct {
	Rule  Rule
	Regex *regexp.Regexp
}

// -----------------------------
// Default Rules
// -----------------------------

var defaultConfig = []Rule{
	{
		Name:   "ssh_success",
		Type:   "success",
		Metric: "ssh_logins",
		Regex:  `Accepted .* for ([^ ]+) from ([0-9.]+)`,
		Labels: []string{"user", "ip"},
	},
	{
		Name:   "ssh_failed",
		Type:   "fail",
		Metric: "ssh_logins",
		Regex:  `Failed .* for ([^ ]+) from ([0-9.]+)`,
		Labels: []string{"user", "ip"},
	},
	{
		Name:   "session_open",
		Type:   "info",
		Metric: "ssh_sessions",
		Regex:  `pam_unix\(sshd:session\): session opened for user ([^ (]+)`,
		Labels: []string{"user"},
	},
	{
		Name:   "session_close",
		Type:   "info",
		Metric: "ssh_sessions",
		Regex:  `pam_unix\(sshd:session\): session closed for user ([^ ]+)`,
		Labels: []string{"user"},
	},
	{
		Name:   "su_open",
		Type:   "info",
		Metric: "ssh_sessions",
		Regex:  `pam_unix\(su:session\): session opened for user ([^ (]+)`,
		Labels: []string{"user"},
	},
	{
		Name:   "su_close",
		Type:   "info",
		Metric: "ssh_sessions",
		Regex:  `pam_unix\(su:session\): session closed for user ([^ ]+)`,
		Labels: []string{"user"},
	},
	{
		Name:   "disconnect",
		Type:   "info",
		Metric: "ssh_events",
		Regex:  `Disconnected from user ([^ ]+)`,
		Labels: []string{"user"},
	},
}

// -----------------------------
// CLI Rule Parsing
// -----------------------------

func ParseRules(input []string) ([]Rule, error) {

	if len(input) == 0 {
		return defaultConfig, nil
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
			Labels: []string{},
		}

		switch {
		case strings.Contains(name, "success"):
			rule.Type = "success"
			rule.Metric = "ssh_logins"
			rule.Labels = []string{"user", "ip"}

		case strings.Contains(name, "fail"):
			rule.Type = "fail"
			rule.Metric = "ssh_logins"
			rule.Labels = []string{"user", "ip"}

		case strings.Contains(name, "session"):
			rule.Type = "info"
			rule.Metric = "ssh_sessions"
			rule.Labels = []string{"user"}

		case strings.Contains(name, "su"):
			rule.Type = "info"
			rule.Metric = "ssh_sessions"
			rule.Labels = []string{"user"}

		case strings.Contains(name, "disconnect"):
			rule.Type = "info"
			rule.Metric = "ssh_events"
			rule.Labels = []string{"user"}
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

// -----------------------------
// Parser
// -----------------------------

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

	return &Parser{rules: compiled}
}

// -----------------------------
// SAFE Parse (NO PANIC EVER)
// -----------------------------

func (p *Parser) Parse(line string) {

	for _, r := range p.rules {

		m := r.Regex.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		// m[0] = full match
		groups := m[1:]

		// safety check
		if len(groups) < len(r.Rule.Labels) {
			continue
		}

		values := groups[:len(r.Rule.Labels)]

		switch r.Rule.Metric {

		case "ssh_logins":
			// status + user + ip REQUIRED
			if len(values) != 2 {
				continue
			}

			status := r.Rule.Type // success/fail
			user := values[0]
			ip := values[1]

			sshLogins.WithLabelValues(status, user, ip).Inc()

		case "ssh_sessions":
			if len(values) == 1 {
				sshSessions.WithLabelValues(r.Rule.Name, values[0]).Inc()
			}

		case "ssh_events":
			if len(values) == 1 {
				sshEvents.WithLabelValues(r.Rule.Type, values[0]).Inc()
			}
		}

		return
	}
}
