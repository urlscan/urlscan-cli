## urlscan scan submit

Submit a URL to scan

```
urlscan scan submit <url> [flags]
```

### Examples

```
  urlscan scan submit <url>
  echo "<url>" | urlscan scan submit -
```

### Options

```
  -c, --country string           Specify which country the scan should be performed from (2-Letter ISO-3166-1 alpha-2 country
  -a, --customagent string      Override User-Agent for this scan
  -h, --help                    help for submit
  -m, --max-wait int            Maximum wait time in seconds (default 60)
  -o, --overrideSafety string   If set to any value, this will disable reclassification of URLs with potential PII in them
  -r, --referer string          Override HTTP referer for this scan
  -t, --tags stringArray        User-defined tags to annotate this scan
  -v, --visibility string       One of public, unlisted, private
  -w, --wait                    Wait for the scan to finish
```

### SEE ALSO

* [urlscan scan](urlscan_scan.md)	 - Scan sub-commands

