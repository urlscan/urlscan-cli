## urlscan pro live-scan store

Store the temporary scan as a permanent snapshot

```
urlscan pro live-scan store [flags]
```

### Examples

```
  urlscan pro live-scan store <scan-id> -S <scanner-id>
  echo <scan-id> | urlscan pro live-scan store - -s <scanner-id>
```

### Options

```
  -h, --help                help for store
  -s, --scanner-id string   ID of the scanner (required)
  -v, --visibility string   Visibility of the scan (public, unlisted or private) (default "private")
```

### SEE ALSO

* [urlscan pro live-scan](urlscan_pro_live-scan.md)	 - Live-scan sub-commands

