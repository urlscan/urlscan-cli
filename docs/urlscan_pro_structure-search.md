## urlscan pro structure-search

Get structurally similar hits to a specific scan

```
urlscan pro structure-search <uuid> [flags]
```

### Examples

```
  urlscan pro structure-search <uuid>
  echo "<uuid>" | urlscan pro structure-search -
```

### Options

```
  -h, --help                  help for structure-search
  -l, --limit int             Maximum number of results that will be returned by the iterator (default 10000)
  -q, --query string          Query to search for in the structure of the scan
      --search-after string   For retrieving the next batch of results, value of the sort attribute of the last (oldest) result you received (comma-separated)
  -s, --size int              Number of results returned by the iterator in each batch (default 100)
```

### SEE ALSO

* [urlscan pro](urlscan_pro.md)	 - Pro sub-commands

