## urlscan pro hostname

Get the historical observations for a specific hostname in the hostname data source

### Synopsis

To have the same idiom with the search command, this command has the following specs:

- Request:
  - limit: the maximum number of results that will be returned by the iterator.
  - size: the number of results returned by the iterator in each batch (equivalent to the API endpoint's "limit" query parameter).
- Response:
  - hasMore: indicates more results are available.

```
urlscan pro hostname [flags]
```

### Examples

```
  urlscan pro hostname <hostname>
  echo "<hostname>" | urlscan pro hostname -
```

### Options

```
  -h, --help                help for hostname
  -l, --limit int           Maximum number of results that will be returned by the iterator (default 10000)
      --all            Return all results; limit is ignored if --all is specified
  -p, --page-state string   Returns additional results starting from this page state from the previous API call
  -s, --size int            Number of results returned by the iterator in each batch (default 100)
```

### SEE ALSO

* [urlscan pro](urlscan_pro.md)	 - Pro sub-commands

