## urlscan pro livescan store

Store the temporary scan as a permanent snapshot

```
urlscan pro livescan store [flags]
```

### Examples

```
  urlscan pro livescan store <scan-id> -S <scanner-id>
  echo <scan-id> | urlscan pro livescan store - -s <scanner-id>
```

### Options

```
  -h, --help                help for store
  -s, --scanner-id string   ID of the scanner (required)
  -v, --visibility string   Visibility of the scan (public, unlisted or private) (default "private")
```

### SEE ALSO

* [urlscan pro livescan](urlscan_pro_livescan.md)	 - Livescan sub-commands

