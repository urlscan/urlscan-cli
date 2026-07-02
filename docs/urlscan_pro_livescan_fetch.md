## urlscan pro livescan fetch

Fetch a response or a download of a live scan by SHA256 file hash

```
urlscan pro livescan fetch [flags]
```

### Examples

```
  urlscan pro livescan fetch <file-hash> -s <scanner-id>
  echo <file-hash> | urlscan pro livescan fetch - -s <scanner-id>
```

### Options

```
  -f, --force               Force overwrite an existing file
  -h, --help                help for fetch
  -o, --output string       Output file name (default <file-hash>)
  -s, --scanner-id string   ID of the scanner (required)
```

### SEE ALSO

* [urlscan pro livescan](urlscan_pro_livescan.md)	 - Livescan sub-commands

