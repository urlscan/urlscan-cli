## urlscan pro livescan response

Get a response of a live scan by SHA256 file hash

```
urlscan pro livescan response [flags]
```

### Examples

```
  urlscan pro livescan response <file-hash> -s <scanner-id>
  echo <file-hash> | urlscan pro livescan response - -s <scanner-id>
```

### Options

```
  -f, --force               Force overwrite an existing file
  -h, --help                help for response
  -o, --output string       Output file name (default <file-hash>)
  -s, --scanner-id string   ID of the scanner (required)
```

### SEE ALSO

* [urlscan pro livescan](urlscan_pro_livescan.md)	 - Livescan sub-commands

