## urlscan scan response

Get a response by SHA256 file hash

```
urlscan scan response <flile-hash> [flags]
```

### Examples

```
  urlscan scan response <file-hash>
  echo "<file-hash>" | urlscan scan response -
```

### Options

```
  -P, --directory-prefix string   Set directory prefix where file will be saved (default ".")
  -f, --force                     Force overwrite an existing file
  -h, --help                      help for response
  -o, --output string             Output file name (default <file-hash>)
```

### SEE ALSO

* [urlscan scan](urlscan_scan.md)	 - Scan sub-commands

