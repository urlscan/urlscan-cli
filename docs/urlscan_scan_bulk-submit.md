## urlscan scan bulk-submit

Bulk submit URLs to scan

### Synopsis

Submit multiple URLs to scan in bulk.

This command allows you to submit a list of URLs for scanning in bulk. You can provide URLs via command line arguments or through a file.
Note that the URLs will be validated before submission, and only valid URLs will be processed.

```
urlscan scan bulk-submit <url>... [flags]
```

### Examples

```
  urlscan scan bulk-submit <url>...
  # submit with a file containing URLs per line, space, or tab
  urlscan scan bulk-submit list_of_urls.txt
  # combine the file input and the URL input
  urlscan scan bulk-submit list_of_urls.txt <url>
```

### Options

```
  -c, --country string          Specify which country the scan should be performed from (2-Letter ISO-3166-1 alpha-2 country
  -a, --customagent string      Override User-Agent for this scan
  -h, --help                    help for bulk-submit
      --max-concurrency int     Maximum number of concurrent requests for batch operation (default 5)
  -m, --max-wait int            Maximum wait time per scan in seconds (default 60)
  -o, --overrideSafety string   If set to any value, this will disable reclassification of URLs with potential PII in them
  -r, --referer string          Override HTTP referer for this scan
  -t, --tags stringArray        User-defined tags to annotate this scan
      --timeout int             Timeout for the batch operation in seconds, 0 means no timeout (default 1800)
  -v, --visibility string       One of public, unlisted, private
  -w, --wait                    Wait for the scan(s) to finish
```

### SEE ALSO

* [urlscan scan](urlscan_scan.md)	 - Scan sub-commands

