## urlscan pro incident create

Create a new incident

```
urlscan pro incident create [flags]
```

### Examples

```
  urlscan pro incident create -o <observable>
```

### Options

```
      --channels strings                    Channels
      --countries strings                   Countries
      --countries-per-interval int          Countries per interval (default 1)
      --expire-after int                    Expire after in seconds (default 0)
      --expire-at string                    Expire at (optional)
  -h, --help                                help for create
      --incident-profile string             Incident profile (optional)
  -o, --observable string                   Observable (hostname, domain, IP or URL) (required)
      --scan-interval int                   Scan interval in seconds (default 0)
      --scan-interval-after-malicious int   Scan interval after malicious in seconds (default 0)
      --scan-interval-after-suspended int   Scan interval after suspended in seconds (default 0)
      --scan-interval-mode string           Scan interval mode (automatic or manual) (default "automatic")
      --stop-delay-inactive int             Stop delay inactive in seconds (default 0)
      --stop-delay-malicious int            Stop delay malicious in seconds (optional) (default 0)
      --stop-delay-suspended int            Stop delay suspended in seconds (optional) (default 0)
      --user-agents strings                 User agents
      --user-agents-per-interval int        User agents per interval (default 1)
      --visibility string                   Visibility (unlisted or private) (default "private")
      --watched-attributes strings          Watched attributes
```

### SEE ALSO

* [urlscan pro incident](urlscan_pro_incident.md)	 - Incident sub-commands

