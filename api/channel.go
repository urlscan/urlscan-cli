package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ChannelOptions struct {
	Channel struct {
		Type           string   `json:"type"`
		Name           string   `json:"name"`
		WebhookURL     string   `json:"webhookURL,omitempty"`
		Frequency      string   `json:"frequency,omitempty"`
		EmailAddresses []string `json:"emailAddresses,omitempty"`
		UTCTime        string   `json:"utcTime,omitempty"`
		IsActive       bool     `json:"isActive,omitempty"`
		IsDefault      bool     `json:"isDefault,omitempty"`
		IgnoreTime     bool     `json:"ignoreTime,omitempty"`
		WeekDays       []string `json:"weekDays,omitempty"`
		Permissions    []string `json:"permissions,omitempty"`
	} `json:"channel"`
}

type ChannelOption func(*ChannelOptions)

func WithChannelType(channelType string) ChannelOption {
	return func(opts *ChannelOptions) {
		opts.Channel.Type = channelType
	}
}

func WithChannelName(name string) ChannelOption {
	return func(opts *ChannelOptions) {
		opts.Channel.Name = name
	}
}

func WithChannelWebhookURL(webhookURL string) ChannelOption {
	return func(opts *ChannelOptions) {
		opts.Channel.WebhookURL = webhookURL
	}
}

func WithChannelFrequency(frequency string) ChannelOption {
	return func(opts *ChannelOptions) {
		opts.Channel.Frequency = frequency
	}
}

func WithChannelEmailAddresses(emailAddresses []string) ChannelOption {
	return func(opts *ChannelOptions) {
		opts.Channel.EmailAddresses = emailAddresses
	}
}

func WithChannelUTCTime(utcTime string) ChannelOption {
	return func(opts *ChannelOptions) {
		opts.Channel.UTCTime = utcTime
	}
}

func WithChannelIsActive(isActive bool) ChannelOption {
	return func(opts *ChannelOptions) {
		opts.Channel.IsActive = isActive
	}
}

func WithChannelIsDefault(isDefault bool) ChannelOption {
	return func(opts *ChannelOptions) {
		opts.Channel.IsDefault = isDefault
	}
}

func WithChannelIgnoreTime(ignoreTime bool) ChannelOption {
	return func(opts *ChannelOptions) {
		opts.Channel.IgnoreTime = ignoreTime
	}
}

func WithChannelWeekDays(weekDays []string) ChannelOption {
	return func(opts *ChannelOptions) {
		opts.Channel.WeekDays = weekDays
	}
}

func WithChannelPermissions(permissions []string) ChannelOption {
	return func(opts *ChannelOptions) {
		opts.Channel.Permissions = permissions
	}
}

func validateChannelOptions(o *ChannelOptions) error {
	if o.Channel.Name == "" {
		return errors.New("name is required")
	}

	switch o.Channel.Type {
	case "webhook":
		if len(o.Channel.WebhookURL) == 0 {
			return errors.New("webhook URL is required")
		}
	case "email":
		if len(o.Channel.EmailAddresses) == 0 {
			return errors.New("email addresses are required")
		}
	}

	return nil
}

func newChannelOptions(opts ...ChannelOption) (*ChannelOptions, error) {
	var o ChannelOptions
	for _, fn := range opts {
		fn(&o)
	}

	err := validateChannelOptions(&o)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (c *Client) CreateChannel(opts ...ChannelOption) (*Response, error) {
	channelOpts, err := newChannelOptions(opts...)
	if err != nil {
		return nil, err
	}

	marshalled, err := json.Marshal(channelOpts)
	if err != nil {
		return nil, err
	}

	return c.NewRequest().SetBodyJSONBytes(marshalled).Post(
		PrefixedPath("/user/channels/"),
	)
}

func (c *Client) UpdateChannel(id string, opts ...ChannelOption) (*Response, error) {
	channelOpts, err := newChannelOptions(opts...)
	if err != nil {
		return nil, err
	}

	marshalled, err := json.Marshal(channelOpts)
	if err != nil {
		return nil, err
	}

	return c.NewRequest().SetBodyJSONBytes(marshalled).Put(
		PrefixedPath(fmt.Sprintf("/user/channels/%s", id)),
	)
}
