## urlscan pro channel update

Update a channel

```
urlscan pro channel update [flags]
```

### Examples

```
  urlscan pro channel update <channel-id> -n <name>
  echo <channel-id> | urlscan pro channel update - -n <name>
```

### Options

```
      --email-addresses strings   Email addresses receiving the notifications (required for type: email)
      --frequency string          Frequency of notifications (live, hourly or daily) (optional)
  -h, --help                      help for update
      --ignore-time               Ignore time constraints (default false)
      --is-active                 Set channel active (default true)
      --is-default                Set channel as default (default false)
  -n, --name string               Channel name (required)
      --permissions strings       Permissions (optional; team:read, team:write)
      --type string               Type of channel (webhook or email) (default "webhook")
      --utc-time string           24 hour UTC time that daily emails are sent at (optional)
      --webhook-url string        Webhook URL (required for type: webhook)
      --week-days strings         Days of the week alerts will be generated (Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday)
```

### SEE ALSO

* [urlscan pro channel](urlscan_pro_channel.md)	 - Channel sub-commands

