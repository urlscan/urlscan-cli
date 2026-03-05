package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNewer(t *testing.T) {
	tests := []struct {
		name    string
		current string
		latest  string
		want    bool
	}{
		{"newer patch", "1.0.0", "1.0.1", true},
		{"newer minor", "1.0.0", "1.1.0", true},
		{"newer major", "1.0.0", "2.0.0", true},
		{"same version", "1.0.0", "1.0.0", false},
		{"older version", "1.1.0", "1.0.0", false},
		{"with v prefix", "v1.0.0", "v1.0.1", true},
		{"mixed prefix", "1.0.0", "v1.0.1", true},
		{"prerelease newer", "1.0.0-rc.1", "1.0.0", true},
		{"calver newer", "2023.01.01", "2023.02.01", true},
		{"calver same", "2023.01.01", "2023.01.01", false},
		{"calver with prerelease", "2023.02.01-rc.1", "2023.02.01", true},
		{"dev (empty)", "", "1.0.0", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsNewer(tt.current, tt.latest))
		})
	}
}
