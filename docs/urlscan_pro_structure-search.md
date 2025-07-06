## urlscan pro structure-search

Get structurally similar results to a specific scan

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
      --all                   Return all results; limit is ignored if --all is specified (default false)
  -h, --help                  help for structure-search
  -l, --limit int             Maximum number of results that will be returned by the iterator (default 10000)
  -q, --query string          Additional query filter
      --search-after string   For retrieving the next batch of results, value of the sort attribute of the last (oldest) result you received (comma-separated)
  -s, --size int              Number of results returned by the iterator in each batch (default 1000)
```

### SEE ALSO

* [urlscan pro](urlscan_pro.md)	 - Pro sub-commands

