## urlscan scan dom

Download a dom by UUID

```
urlscan scan dom <uuid> [flags]
```

### Examples

```
  urlscan scan dom <uuid>
  echo "<uuid>" | urlscan scan dom -
```

### Options

```
  -P, --directory-prefix string   Set directory prefix where file will be saved (default ".")
  -f, --force                     Force overwrite an existing file
  -h, --help                      help for dom
  -o, --output string             Output file name (default <uuid>.html)
```

### SEE ALSO

* [urlscan scan](urlscan_scan.md)	 - Scan sub-commands

