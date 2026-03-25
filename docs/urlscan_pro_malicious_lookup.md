## urlscan pro malicious lookup

Look up how often an observable has been seen in malicious scan results

### Synopsis

Look up how often an observable has been seen in malicious scan results, along with first and last seen timestamps. Type must be one of: ip, hostname, domain, url.

```
urlscan pro malicious lookup <type> <value> [flags]
```

### Examples

```
  urlscan pro malicious lookup ip 192.0.2.1
  urlscan pro malicious lookup hostname www.example.com
  urlscan pro malicious lookup domain example.com
  urlscan pro malicious lookup url "https://example.com/path"
  echo "192.0.2.1" | urlscan pro malicious lookup ip -
```

### Options

```
  -h, --help     help for lookup
      --refang   Refang an input (convert '[.]' back to '.' and so on)
```

### SEE ALSO

* [urlscan pro malicious](urlscan_pro_malicious.md)	 - Malicious sub-commands

