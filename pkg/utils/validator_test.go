package utils

import (
	"testing"
)

func TestValidateULID(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		// Valid ULID
		{"01ARZ3NDEKTSV4RRFFQ69G5FAV", false},
		// Invalid length
		{"01ARZ3NDEKTSV4RRFFQ69G5FA", true},
		{"01ARZ3NDEKTSV4RRFFQ69G5FAVV", true},
		// Contains invalid characters (lowercase)
		{"01arz3ndektsv4rrffq69g5fav", true},
		// Contains invalid characters (I, L, O, U)
		{"01ARZ3NDEKTSV4RRFFQ69G5FAI", true},
		{"01ARZ3NDEKTSV4RRFFQ69G5FAL", true},
		{"01ARZ3NDEKTSV4RRFFQ69G5FAO", true},
		{"01ARZ3NDEKTSV4RRFFQ69G5FAU", true},
		// Contains invalid characters (symbols)
		{"01ARZ3NDEKTSV4RRFFQ69G5FA!", true},
		// All digits
		{"01234567890123456789012345", false},
		// All valid uppercase letters
		{"ABCDEFGHJKMNPQRSTVWXYZ2345", false},
	}

	for _, tt := range tests {
		err := ValidateULID(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateULID(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		// Valid UUIDs
		{"123e4567-e89b-12d3-a456-426614174000", false},
		{"550e8400-e29b-41d4-a716-446655440000", false},
		{"FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF", false},
		{"00000000-0000-0000-0000-000000000000", false},
		// Invalid length
		{"123e4567-e89b-12d3-a456-42661417400", true},
		{"123e4567-e89b-12d3-a456-4266141740000", true},
		// Missing hyphens
		{"123e4567e89b12d3a456426614174000", true},
		// Hyphens in wrong places
		{"123e4567e-89b1-2d3a-4564-26614174000", true},
		// Invalid characters (symbols)
		{"123e4567-e89b-12d3-a456-42661417400!", true},
		// Invalid characters (non-hex)
		{"123e4567-e89b-12d3-a456-42661417400g", true},
		// Uppercase valid
		{"ABCDEFAB-CDEF-ABCD-EFAB-CDEFABCDEFAB", false},
		// Lowercase valid
		{"abcdefab-cdef-abcd-efab-cdefabcdefab", false},
	}

	for _, tt := range tests {
		err := ValidateUUID(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateUUID(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		// Valid URLs
		{"http://example.com", false},
		{"https://example.com", false},
		{"http://example.com/path?query=1", false},
		{"https://sub.domain.com:8080/path/to/resource", false},
		{"http://example.com/#fragment", false},
		{"https://example.com/?q=test#frag", false},
		{"https://1.1.1.1", false},
		// Invalid URLs
		{"example.com", true},
		{"htp://example.com", true},
		{"http:/example.com", true},
		{"http//example.com", true},
		{"http://", true},
		{"http:// example.com", true},
		{"http://example .com", true},
		{"", true},
		{"https://", true},
		{"ftp:/example.com", true},
		{"http://exa mple.com", true},
		{"invalid-input", true},
	}

	for _, tt := range tests {
		err := ValidateURL(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateURL(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

func TestValidateDomain(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		// Valid domains
		{"example.com", false},
		{"sub.domain.com", false},
		{"example.co.uk", false},
		{"my-site.org", false},
		// Invalid domains
		{"example .com", true},
		{" example.com", true},
		{"example.com ", true},
		{"exa mple.com", true},
		{"", true},
		{".com", true},
		{"com.", true},
		{"invalid-input", true},
		{"192.168.1", true},
	}

	for _, tt := range tests {
		err := ValidateDomain(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateDomain(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

func TestValidateIP(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		// Valid IPv4
		{"192.168.1.1", false},
		{"255.255.255.255", false},
		{"0.0.0.0", false},
		// Invalid IPv4
		{"192.168.1.256", true},
		{"192.168.1", true},
		{"192.168.1.1.1", true},
		{"192.168.1.-1", true},
		{"192.168.1.1a", true},
		{"", true},
		{"invalid-ip", true},
		// Valid IPv6
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", false},
		{"2001:db8:85a3::8a2e:370:7334", false},
		{"::1", false},
		// Invalid IPv6
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334:1234", true},
		{"2001:db8:85a3::8a2e::7334", true},
		{"2001:db8:85a3::8a2e:370g:7334", true},
	}

	for _, tt := range tests {
		err := ValidateIP(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateIP(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

func TestValidateNetworkIndicator(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		// Valid IPs
		{"192.168.1.1", false},
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", false},
		// Valid Domains
		{"example.com", false},
		{"sub.domain.com", false},
		// Valid URLs
		{"http://example.com", false},
		{"https://example.com/path?query=1", false},
		{"https://sub.domain.com:8080/path/to/resource", false},
		{"http://example.com/#fragment", false},
		{"https://example.com/?q=test#frag", false},
		// Invalid inputs
		{"192.168.1", true},
		{"example .com", true},
		{"", true},
		{"invalid-input", true},
	}

	for _, tt := range tests {
		err := ValidateNetworkIndicator(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateNetworkIndicator(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}
