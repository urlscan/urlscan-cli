## urlscan pro livescan download

Download a resource of a live scan by SHA256 file hash

```
urlscan pro livescan download [flags]
```

### Examples

```
  urlscan pro livescan download <file-hash> -s <scanner-id>
  echo <file-hash> | urlscan pro livescan download - -s <scanner-id>
```

### Options

```
  -f, --force               Force overwrite an existing file
  -h, --help                help for download
  -o, --output string       Output file name (default <file-hash>)
  -s, --scanner-id string   ID of the scanner (required)
```

### SEE ALSO

* [urlscan pro livescan](urlscan_pro_livescan.md)	 - Livescan sub-commands

