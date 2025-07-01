## urlscan pro subscription create

Create a new subscription

```
urlscan pro subscription create [flags]
```

### Examples

```
  urlscan pro subscription create -s <search-id-1> -s <search-id-2> -f <frequency> -e <email-address-1> -e <email-address-2> -n <name>
```

### Options

```
  -d, --description string        Description of the subscription (optional)
  -e, --email-addresses strings   Email addresses to send notifications to (required)
  -f, --frequency string          Frequency of notifications (live, hourly or daily) (required)
  -h, --help                      help for create
  -t, --ignore-time               Whether to ignore time constraints (required) (default false)
  -a, --is-active                 Whether the subscription is active (required) (default true)
  -n, --name string               Name of the subscription (required)
  -s, --search-ids strings        Array of search IDs associated with this subscription (required)
  -i, --subscription-id string    Subscription ID (optional, if not provided a new id will be generated)
```

### SEE ALSO

* [urlscan pro subscription](urlscan_pro_subscription.md)	 - Subscription sub-commands

