## urlscan hostname

Get the historical observations for a specific hostname in the hostname data source

```
urlscan hostname [flags]
```

### Examples

```
  urlscan hostname <hostname>
  echo "<hostname>" | urlscan hostname -
```

### Options

```
  -h, --help                help for hostname
  -l, --limit int           Return at most this many results (Minimum 10, Maximum 10,000) (default 1000)
  -p, --page-state string   Continue return additional results starting from this page state from the previous API call
```

### SEE ALSO

* [urlscan](urlscan.md)	 - A CLI tool for interacting with urlscan.io

