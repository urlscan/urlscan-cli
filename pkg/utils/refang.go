package utils

import "strings"

// Refang reverses common defanging transformations as described in https://www.ietf.org/archive/id/draft-grimminck-safe-ioc-sharing-00.html
func Refang(s string) string {
	replacer := strings.NewReplacer(
		"hxxps://", "https://",
		"hXXps://", "https://",
		"hXXPs://", "https://",
		"HXXPS://", "HTTPS://",
		"hxxp://", "http://",
		"hXXp://", "http://",
		"hXXP://", "http://",
		"HXXP://", "HTTP://",
		"[.]", ".",
		"[@]", "@",
		"[:]", ":",
	)
	return replacer.Replace(s)
}
