package utils

import (
	"fmt"
	"regexp"
)

func ValidateULID(s string) error {
	// A ULID is a 26-character string using Crockford's Base32 (0-9, A-Z except I, L, O, U)
	if len(s) != 26 {
		return fmt.Errorf("invalid ULID format: %s", s)
	}
	re := regexp.MustCompile(`^[0-9A-HJKMNPQRSTVWXYZ]{26}$`)
	if !re.Match([]byte(s)) {
		return fmt.Errorf("invalid ULID format: %s", s)
	}
	return nil
}

func ValidateUUID(s string) error {
	// A UUID is a 36-character string in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if len(s) != 36 {
		return fmt.Errorf("invalid UUID format: %s", s)
	}
	re := regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`)
	if !re.Match([]byte(s)) {
		return fmt.Errorf("invalid UUID format: %s", s)
	}
	return nil
}
