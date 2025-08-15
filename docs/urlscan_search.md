## urlscan search

Search by a query

```
urlscan search <query> [flags]
```

### Examples

```
  urlscan search <query>
  echo "<query>" | urlscan search -
```

### Options

```
      --all                   Return all results; limit is ignored if --all is specified (default false)
  -D, --datasource string     Datasources to search: scans (urlscan.io), hostnames, incidents, notifications, certificates (urlscan Pro) (default "scans")
  -h, --help                  help for search
  -l, --limit int             Maximum number of results that will be returned by the iterator (default 10000)
      --search-after string   For retrieving the next batch of results, value of the sort attribute of the last (oldest) result you received (comma-separated)
  -s, --size int              Number of results returned by the iterator in each batch (default 100)
```

### SEE ALSO

* [urlscan](urlscan.md)	 - A CLI tool for interacting with urlscan.io

