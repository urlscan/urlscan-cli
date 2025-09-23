package api

import (
	"testing"
)

func TestNewIncidentOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    []IncidentOption
		wantErr bool
	}{
		{
			name:    "valid incident options",
			opts:    []IncidentOption{WithIncidentObservable("test.com")},
			wantErr: false,
		},
		{
			name:    "missing observable",
			opts:    []IncidentOption{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newIncidentOptions(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("newIncidentOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
