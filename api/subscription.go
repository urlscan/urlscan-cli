package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

type SubscriptionOptions struct {
	Subscription struct {
		SearchIds            []string `json:"searchIds"`
		Frequency            string   `json:"frequency"`
		EmailAddresses       []string `json:"emailAddresses"`
		Name                 string   `json:"name"`
		Description          string   `json:"description,omitempty"`
		IsActive             bool     `json:"isActive"`
		IgnoreTime           bool     `json:"ignoreTime"`
		WeekDays             []string `json:"weekDays,omitempty"`
		Permissions          []string `json:"permissions,omitempty"`
		ChannelIds           []string `json:"channelIds,omitempty"`
		IncidentChannelIds   []string `json:"incidentChannelIds,omitempty"`
		IncidentProfileId    string   `json:"incidentProfileId,omitempty"`
		IncidentVisibility   string   `json:"incidentVisibility,omitempty"`
		IncidentCreationMode string   `json:"incidentCreationMode,omitempty"`
		IncidentWatchKeys    string   `json:"incidentWatchKeys,omitempty"`
	} `json:"subscription"`
}

type SubscriptionOption func(*SubscriptionOptions)

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

func WithSubscriptionWeekDays(weekDays []string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.WeekDays = weekDays
	}
}

func WithSubscriptionPermissions(permissions []string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.Permissions = permissions
	}
}

func WithSubscriptionChannelIds(channelIds []string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.ChannelIds = channelIds
	}
}

func WithSubscriptionIncidentChannelIds(incidentChannelIds []string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.IncidentChannelIds = incidentChannelIds
	}
}

func WithSubscriptionIncidentProfileId(incidentProfileId string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.IncidentProfileId = incidentProfileId
	}
}

func WithSubscriptionIncidentVisibility(incidentVisibility string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.IncidentVisibility = incidentVisibility
	}
}

func WithSubscriptionIncidentCreationMode(incidentCreationMode string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.IncidentCreationMode = incidentCreationMode
	}
}

func WithSubscriptionIncidentWatchKeys(incidentWatchKeys string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Subscription.IncidentWatchKeys = incidentWatchKeys
	}
}

func validateSubscriptionOptions(o *SubscriptionOptions) error {
	if len(o.Subscription.SearchIds) == 0 {
		return errors.New("search IDs are required")
	}
	if o.Subscription.Frequency == "" {
		return errors.New("frequency is required")
	}
	if len(o.Subscription.EmailAddresses) == 0 {
		return errors.New("email addresses are required")
	}
	if o.Subscription.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

func newSubscriptionOptions(opts ...SubscriptionOption) (*SubscriptionOptions, error) {
	var o SubscriptionOptions
	for _, fn := range opts {
		fn(&o)
	}

	err := validateSubscriptionOptions(&o)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (c *Client) CreateSubscription(opts ...SubscriptionOption) (*Response, error) {
	subscriptionOptions, err := newSubscriptionOptions(opts...)
	if err != nil {
		return nil, err
	}

	marshalled, err := json.Marshal(subscriptionOptions)
	if err != nil {
		return nil, err
	}

	return c.NewRequest().SetBodyJSONBytes(marshalled).Post(
		PrefixedPath("/user/subscriptions/"),
	)
}

func (c *Client) UpdateSubscription(id string, opts ...SubscriptionOption) (*Response, error) {
	subscriptionOptions, err := newSubscriptionOptions(opts...)
	if err != nil {
		return nil, err
	}

	marshalled, err := json.Marshal(subscriptionOptions)
	if err != nil {
		return nil, err
	}

	return c.NewRequest().SetBodyJSONBytes(marshalled).Put(
		PrefixedPath(fmt.Sprintf("/user/subscriptions/%s/", id)),
	)
}
