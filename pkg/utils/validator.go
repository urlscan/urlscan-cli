package utils

import "fmt"

func ValidateULID(s string) error {
	// A ULID is a 26-character string using Crockford's Base32 (0-9, A-Z except I, L, O, U)
	if len(s) != 26 {
		return fmt.Errorf("invalid ULID format: %s", s)
	}
	for _, c := range s {
		switch {
		case c >= '0' && c <= '9':
			// valid
		case c >= 'A' && c <= 'Z' && c != 'I' && c != 'L' && c != 'O' && c != 'U':
			// valid
		default:
			return fmt.Errorf("invalid ULID format: %s", s)
		}
	}
	return nil
}

func ValidateUUID(s string) error {
	// A UUID is a 36-character string in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if len(s) != 36 {
		return fmt.Errorf("invalid UUID format: %s", s)
	}
	for i, c := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return fmt.Errorf("invalid UUID format: %s", s)
			}
		} else if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return fmt.Errorf("invalid UUID format: %s", s)
		}
	}
	return nil
}
