package api

import (
	"testing"
)

func TestNewSubscriptionOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    []SubscriptionOption
		wantErr bool
	}{
		{
			name: "valid subscription options",
			opts: []SubscriptionOption{
				WithSubscriptionSearchIds([]string{"search1"}),
				WithSubscriptionFrequency("daily"),
				WithSubscriptionEmailAddresses([]string{"test@example.com"}),
				WithSubscriptionName("Test Subscription"),
				WithSubscriptionDescription("desc"),
				WithSubscriptionIsActive(true),
				WithSubscriptionIgnoreTime(false),
			},
			wantErr: false,
		},
		{
			name: "missing searchIds",
			opts: []SubscriptionOption{
				WithSubscriptionID("subid"),
				WithSubscriptionSearchIds([]string{}),
				WithSubscriptionFrequency("daily"),
				WithSubscriptionEmailAddresses([]string{"test@example.com"}),
				WithSubscriptionName("Test Subscription"),
			},
			wantErr: true,
		},
		{
			name: "missing frequency",
			opts: []SubscriptionOption{
				WithSubscriptionID("subid"),
				WithSubscriptionSearchIds([]string{"search1"}),
				WithSubscriptionFrequency(""),
				WithSubscriptionEmailAddresses([]string{"test@example.com"}),
				WithSubscriptionName("Test Subscription"),
			},
			wantErr: true,
		},
		{
			name: "missing emailAddresses",
			opts: []SubscriptionOption{
				WithSubscriptionID("subid"),
				WithSubscriptionSearchIds([]string{"search1"}),
				WithSubscriptionFrequency("daily"),
				WithSubscriptionEmailAddresses([]string{}),
				WithSubscriptionName("Test Subscription"),
			},
			wantErr: true,
		},
		{
			name: "missing name",
			opts: []SubscriptionOption{
				WithSubscriptionID("subid"),
				WithSubscriptionSearchIds([]string{"search1"}),
				WithSubscriptionFrequency("daily"),
				WithSubscriptionEmailAddresses([]string{"test@example.com"}),
				WithSubscriptionName(""),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newSubscriptionOptions(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("newSubscriptionOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
