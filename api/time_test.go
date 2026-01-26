package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetLast7Days(t *testing.T) {
	expected := []string{
		"20260101",
		"20260102",
		"20260103",
		"20260104",
		"20260105",
		"20260106",
		"20260107",
	}
	assert.Equal(t, expected, GetLast7Days(time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC)))
}
