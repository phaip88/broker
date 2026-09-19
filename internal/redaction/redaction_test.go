package redaction

import "testing"

func TestRedactRulesAndHeaders(t *testing.T) {
	got, tokens := Redact("Authorization: Bearer abc123\nCookie: sid=xyz", map[string]string{"abc123": "TOKEN_1"})
	if got == "" || tokens["TOKEN_1"] == "" { t.Fatal("expected tokenized output") }
	if got == "Authorization: Bearer abc123\nCookie: sid=xyz" { t.Fatal("sensitive values were not redacted") }
}
