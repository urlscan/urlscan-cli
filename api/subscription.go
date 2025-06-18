package api

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
	subscriptionOptions := &SubscriptionOptions{}

	for _, opt := range opts {
		opt(subscriptionOptions)
	}

	return subscriptionOptions
}
