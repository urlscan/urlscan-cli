## urlscan pro subscription update

Update a subscription

```
urlscan pro subscription update [flags]
```

### Examples

```
  urlscan pro subscription update <subscription-id> -s <search-id-1> -s <search-id-2> -f <frequency> -e <email-address-1> -e <email-address-2> -n <name>
```

### Options

```
  -d, --description string        Description of the subscription (optional)
  -e, --email-addresses strings   Email addresses to send notifications to (required)
  -f, --frequency string          Frequency of notifications (required/live, hourly or daily)
  -h, --help                      help for update
  -t, --ignore-time               Whether to ignore time constraints (required/defaults to false)
  -a, --is-active                 Whether the subscription is active (required/defaults to true) (default true)
  -n, --name string               Name of the subscription (required)
  -s, --search-ids strings        Array of search IDs associated with this subscription (required)
```

### SEE ALSO

* [urlscan pro subscription](urlscan_pro_subscription.md)	 - Subscription sub-commands

