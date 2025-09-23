package api

import (
	"testing"
)

func TestNewChannelOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    []ChannelOption
		wantErr bool
	}{
		{
			name:    "missing name",
			opts:    []ChannelOption{WithChannelType("webhook")},
			wantErr: true,
		},
		{
			name:    "missing webhook url",
			opts:    []ChannelOption{WithChannelType("webhook"), WithChannelName("test")},
			wantErr: true,
		},
		{
			name:    "with webhook url",
			opts:    []ChannelOption{WithChannelType("webhook"), WithChannelName("test"), WithChannelWebhookURL("https://example.com")},
			wantErr: false,
		},
		{
			name:    "missing email addresses",
			opts:    []ChannelOption{WithChannelType("email"), WithChannelName("test")},
			wantErr: true,
		},
		{
			name:    "with email addresses",
			opts:    []ChannelOption{WithChannelType("email"), WithChannelName("test"), WithChannelEmailAddresses([]string{"foo@example.com"})},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newChannelOptions(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("newChannelOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
