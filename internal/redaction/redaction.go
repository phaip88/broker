package redaction

import (
	"regexp"
	"strings"
)

// Redact performs deterministic, typed masking for a first local prototype.
// Production replacement: project-scoped encrypted token map + policy engine.
func Redact(text string, rules map[string]string) (string, map[string]string) {
	tokens := map[string]string{}
	for raw, token := range rules {
		if raw == "" || token == "" { continue }
		text = strings.ReplaceAll(text, raw, token)
		tokens[token] = "redacted"
	}
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(authorization\s*[:=]\s*bearer\s+)[^\s,]+`),
		regexp.MustCompile(`(?i)(cookie\s*[:=]\s*)[^\r\n]+`),
	}
	text = patterns[0].ReplaceAllString(text, `${1}TOKEN_REDACTED`)
	text = patterns[1].ReplaceAllString(text, `${1}COOKIE_REDACTED`)
	return text, tokens
}
