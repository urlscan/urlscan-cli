package utils

import "testing"

func TestRefang(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Protocol replacements
		{"hxxps://evil[.]example[.]com/path", "https://evil.example.com/path"},
		{"hxxp://evil[.]example[.]com/path", "http://evil.example.com/path"},
		// Case variations
		{"hXXps://evil[.]example[.]com", "https://evil.example.com"},
		{"HXXPS://EVIL[.]EXAMPLE[.]COM", "HTTPS://EVIL.EXAMPLE.COM"},
		{"hXXp://evil[.]example[.]com", "http://evil.example.com"},
		{"HXXP://EVIL[.]EXAMPLE[.]COM", "HTTP://EVIL.EXAMPLE.COM"},
		// FTP
		{"fxp://files[.]example[.]com", "ftp://files.example.com"},
		{"FXP://FILES[.]EXAMPLE[.]COM", "FTP://FILES.EXAMPLE.COM"},
		// Email addresses
		{"user[@]phishing[.]example[.]com", "user@phishing.example.com"},
		// Credentials in URI
		{"hxxp://username:password[@]attacker[.]com", "http://username:password@attacker.com"},
		// IPv6 with colon defanging
		{"hxxp://[2001:db8[:]1]:8080", "http://[2001:db8:1]:8080"},
		// Domain only
		{"evil[.]example[.]com", "evil.example.com"},
		// No defanging present (passthrough)
		{"https://example.com", "https://example.com"},
		{"", ""},
		// Multiple dots
		{"a[.]b[.]c[.]d[.]e", "a.b.c.d.e"},
		// Mixed: some defanged, some not
		{"hxxps://example.com/foo[.]bar", "https://example.com/foo.bar"},
	}

	for _, tt := range tests {
		got := Refang(tt.input)
		if got != tt.want {
			t.Errorf("Refang(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
