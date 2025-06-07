## urlscan pro hostname

Get the historical observations for a specific hostname in the hostname data source

### Synopsis

To have the same idiom with the search command, this command has the following specs:

- Request:
  - limit: is not exactly same as the API endpoint's "limit" query parameter. It is the maximum number of results that will be returned by the iterator.
  - size: is equivalent to the API endpoint's "limit" query parameter, which is the number of results returned by the iterator in each batch.
- Response:
  - hasMore: is an additional field (not included in the API endpoint response) indicates if there are more results available.


```
urlscan pro hostname [flags]
```

### Examples

```
  urlscan hostname <hostname>
  echo "<hostname>" | urlscan hostname -
```

### Options

```
  -h, --help                help for hostname
  -l, --limit int           Maximum number of results that will be returned by the iterator (default 10000)
  -p, --page-state string   Continue return additional results starting from this page state from the previous API call
  -s, --size int            Number of results returned by the iterator in each batch (default 1000)
```

### SEE ALSO

* [urlscan pro](urlscan_pro.md)	 - Pro sub-commands

