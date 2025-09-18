package utils

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"slices"
	"strings"
)

var (
	re_ULID       = regexp.MustCompile(`^[0-9A-HJKMNPQRSTVWXYZ]{26}$`)
	re_UUID       = regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`)
	re_dot        = regexp.MustCompile(`\.`)
	re_digit_only = regexp.MustCompile(`^\d+$`)
)

func ValidateULID(s string) error {
	// A ULID is a 26-character string using Crockford's Base32 (0-9, A-Z except I, L, O, U)
	if len(s) != 26 || !re_ULID.Match([]byte(s)) {
		return fmt.Errorf("invalid ULID format: %s", s)
	}
	return nil
}

func ValidateUUID(s string) error {
	// A UUID is a 36-character string in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if len(s) != 36 || !re_UUID.Match([]byte(s)) {
		return fmt.Errorf("invalid UUID format: %s", s)
	}
	return nil
}

func ValidateURL(s string) error {
	parsed, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid URL format: %s", s)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("invalid URL format: %s", s)
	}

	if parsed.Host == "" {
		return fmt.Errorf("invalid URL format: %s", s)
	}

	// host should be a valid domain or an IP address
	if ValidateDomain(parsed.Host) != nil && ValidateIP(parsed.Host) != nil {
		return fmt.Errorf("invalid URL format: %s", s)
	}

	return nil
}

func ValidateDomain(s string) error {
	// check whether it starts with https:// or http://
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return fmt.Errorf("invalid domain format: %s", s)
	}

	parsed, err := url.Parse("http://" + s)
	if err != nil {
		return fmt.Errorf("invalid domain format: %s", s)
	}

	if parsed.Host == "" {
		return fmt.Errorf("invalid domain format: %s", s)
	}

	parts := re_dot.Split(parsed.Host, -1)
	// should have one or more dots
	if len(parts) <= 1 {
		return fmt.Errorf("invalid domain format: %s", s)
	}
	// ensure that there is no empty part in the domain (e.g. example..com)
	if slices.Contains(parts, "") {
		return fmt.Errorf("invalid domain format: %s", s)
	}

	// should not all the parts are digits to avoid IP addresses like string (e.g. 192.168.1)
	allDigits := true
	for _, part := range parts {
		if !re_digit_only.MatchString(part) {
			allDigits = false
			break
		}
	}
	if allDigits {
		return fmt.Errorf("invalid domain format: %s", s)
	}

	return nil
}

func ValidateIP(s string) error {
	if net.ParseIP(s) == nil {
		return fmt.Errorf("invalid IP format: %s", s)
	}
	return nil
}

func ValidateNetworkIndicator(s string) error {
	if ValidateIP(s) == nil {
		return nil
	}

	if ValidateDomain(s) == nil {
		return nil
	}

	if ValidateURL(s) == nil {
		return nil
	}

	return fmt.Errorf("invalid network indicator format: %s", s)
}
