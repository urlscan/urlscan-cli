package api

import (
	"encoding/json"
	"fmt"
)

type SubscriptionOptions struct {
	Subscription struct {
		ID             string   `json:"_id"`
		SearchIds      []string `json:"searchIds"`
		Frequency      string   `json:"frequency"`
		EmailAddresses []string `json:"emailAddresses"`
		Name           string   `json:"name"`
		Description    string   `json:"description,omitempty"`
		IsActive       bool     `json:"isActive"`
		IgnoreTime     bool     `json:"ignoreTime"`
	} `json:"subscription"`
}

type SubscriptionOption func(*SubscriptionOptions)

func WithSubscriptionID(id string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.ID = id
	}
}

func WithSubscriptionSearchIds(searchIds []string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.SearchIds = searchIds
	}
}

func WithSubscriptionFrequency(frequency string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.Frequency = frequency
	}
}

func WithSubscriptionEmailAddresses(emailAddresses []string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.EmailAddresses = emailAddresses
	}
}

func WithSubscriptionName(name string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.Name = name
	}
}

func WithSubscriptionDescription(description string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.Description = description
	}
}

func WithSubscriptionIsActive(isActive bool) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.IsActive = isActive
	}
}

func WithSubscriptionIgnoreTime(ignoreTime bool) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.IgnoreTime = ignoreTime
	}
}

func newSubscriptionOptions(opts ...SubscriptionOption) *SubscriptionOptions {
	var o SubscriptionOptions
	for _, fn := range opts {
		fn(&o)
	}
	return &o
}

func (c *Client) CreateSubscription(opts ...SubscriptionOption) (*Response, error) {
	subscriptionOptions := newSubscriptionOptions(opts...)
	marshalled, err := json.Marshal(subscriptionOptions)
	if err != nil {
		return nil, err
	}

	return c.NewRequest().SetBodyJSONBytes(marshalled).Post(
		PrefixedPath("/user/subscriptions/"),
	)
}

func (c *Client) UpdateSubscription(opts ...SubscriptionOption) (*Response, error) {
	subscriptionOptions := newSubscriptionOptions(opts...)
	marshalled, err := json.Marshal(subscriptionOptions)
	if err != nil {
		return nil, err
	}

	return c.NewRequest().SetBodyJSONBytes(marshalled).Put(
		PrefixedPath(fmt.Sprintf("/user/subscriptions/%s/", subscriptionOptions.Subscription.ID)),
	)
}
