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
      --channel-ids strings             Array of channel IDs associated with this subscription
  -d, --description string              Description of the subscription (optional)
  -e, --email-addresses strings         Email addresses to send notifications to (required)
  -f, --frequency string                Frequency of notifications (live, hourly or daily) (required)
  -h, --help                            help for create
  -t, --ignore-time                     Whether to ignore time constraints (default false)
      --incident-channel-ids strings    Array of incident channel IDs associated with this subscription
      --incident-creation-mode string   Incident creation rule (none, default, always, ignore-if-exists)
      --incident-profile-id string      Incident Profile ID associated with this subscription
      --incident-visibility string      Incident visibility (unlisted, private)
      --incident-watch-keys string      Source/key to watch in the incident (scans/page.url, scans/page.domain, scans/page.ip, scans/page.apexDomain, hostnames/hostname, hostnames/domain)
  -a, --is-active                       Whether the subscription is active (default true)
  -n, --name string                     Name of the subscription (required)
      --permissions strings             Permissions (team:read, team:write)
  -s, --search-ids strings              Array of search IDs associated with this subscription (required)
      --week-days strings               Days of the week alerts will be generated (Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday)
```

### SEE ALSO

* [urlscan pro subscription](urlscan_pro_subscription.md)	 - Subscription sub-commands

