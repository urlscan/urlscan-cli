package api

import (
	"encoding/json"
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

func TestChannelOptionsMarshalJSON_WithExtra(t *testing.T) {
	opts := ChannelOptions{} // nolint:exhaustruct
	opts.Extra = map[string]any{
		"customField": "customValue",
	}
	opts.Channel.Type = "webhook"
	opts.Channel.Name = "test"
	opts.Channel.WebhookURL = "https://example.com"

	b, err := json.Marshal(opts)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	// channel field should be present
	if _, ok := m["channel"]; !ok {
		t.Error("expected 'channel' key in JSON output")
	}

	// extra field should be merged at top level
	if v, ok := m["customField"]; !ok || v != "customValue" {
		t.Errorf("expected 'customField' = 'customValue', got %v", v)
	}
}

func TestChannelOptionsMarshalJSON_WithoutExtra(t *testing.T) {
	opts := ChannelOptions{} // nolint:exhaustruct
	opts.Channel.Type = "webhook"
	opts.Channel.Name = "test"
	opts.Channel.WebhookURL = "https://example.com"

	b, err := json.Marshal(opts)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	// should only have "channel" key
	if len(m) != 1 {
		t.Errorf("expected 1 key, got %d: %v", len(m), m)
	}
}
