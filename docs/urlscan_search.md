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
  -h, --help                  help for search
  -l, --limit int             Maximum number of results that will be returned by the iterator (default 10000)
      --no-limit              Don't limit the number of results returned by the iterator, limit is ignored if it's set (default false)
      --search-after string   For retrieving the next batch of results, value of the sort attribute of the last (oldest) result you received (comma-separated)
  -s, --size int              Number of results returned by the iterator in each batch (default 100)
```

### SEE ALSO

* [urlscan](urlscan.md)	 - A CLI tool for interacting with urlscan.io

