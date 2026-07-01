## urlscan pro livescan screenshot

Get a screenshot of a live scan

```
urlscan pro livescan screenshot [flags]
```

### Examples

```
  urlscan pro livescan screenshot <scan-id> -s <scanner-id>
  echo <scan-id> | urlscan pro livescan screenshot - -s <scanner-id>
```

### Options

```
  -f, --force               Force overwrite an existing file
  -h, --help                help for screenshot
  -o, --output string       Output file name (default <uuid>.png)
  -s, --scanner-id string   ID of the scanner (required)
```

### SEE ALSO

* [urlscan pro livescan](urlscan_pro_livescan.md)	 - Livescan sub-commands

