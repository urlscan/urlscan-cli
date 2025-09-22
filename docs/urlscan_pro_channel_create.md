## urlscan pro channel create

Create a new channel

```
urlscan pro channel create [flags]
```

### Examples

```
  urlscan pro channel create -n <name>
```

### Options

```
      --email-addresses strings   Email addresses receiving the notifications
      --frequency string          Frequency of notifications (live, hourly or daily) (optional)
  -h, --help                      help for create
      --ignore-time               Whether to ignore time constraints (default false)
      --is-active                 Whether the channel is active (default true)
      --is-default                Whether the channel is the default channel (default false)
  -n, --name string               Channel name (required)
      --permissions strings       Permissions
      --type string               Type of channel (webhook or email) (default "webhook")
      --utc-time string           24 hour UTC time that daily emails are sent (optional)
      --webhook-url string        Webhook URL (required for type: webhook)
      --week-days strings         Days of the week alerts will be generated (Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday)
```

### SEE ALSO

* [urlscan pro channel](urlscan_pro_channel.md)	 - Channel sub-commands

