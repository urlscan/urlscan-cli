## urlscan pro live-scan dom

Get a dom of a live scan

```
urlscan pro live-scan dom [flags]
```

### Examples

```
  urlscan pro live-scan dom <scan-id> -s <scanner-id>
  echo <scan-id> | urlscan pro live-scan dom - -s <scanner-id>
```

### Options

```
  -f, --force               Force overwrite an existing file.
  -h, --help                help for dom
  -o, --output string       Output file name (default <uuid>.html)
  -s, --scanner-id string   ID of the scanner (required)
```

### SEE ALSO

* [urlscan pro live-scan](urlscan_pro_live-scan.md)	 - Live-scan sub-commands

