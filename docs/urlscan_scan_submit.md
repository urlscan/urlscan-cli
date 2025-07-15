## urlscan scan submit

Submit a URL to scan

```
urlscan scan submit <url> [flags]
```

### Examples

```
  urlscan scan submit <url>...
  echo "<url>" | urlscan scan submit -
```

### Options

```
  -c, --country string          Specify which country the scan should be performed from (2-Letter ISO-3166-1 alpha-2 country
  -a, --customagent string      Override User-Agent for this scan
  -f, --force                   Force overwrite an existing file
  -h, --help                    help for submit
  -m, --max-wait int            Maximum wait time per scan in seconds (default 60)
  -o, --overrideSafety string   If set to any value, this will disable reclassification of URLs with potential PII in them
  -r, --referer string          Override HTTP referer for this scan
  -t, --tags stringArray        User-defined tags to annotate this scan
  -v, --visibility string       One of public, unlisted, private
  -w, --wait                    Wait for the scan(s) to finish
      --with-both               Download both a screenshot and a DOM (this flag overrides wait, with-screen and with-both flags to true)
      --with-dom                Download a DOM (this flag overrides wait flag to true)
      --with-screenshot         Download a screenshot (this flag overrides wait flag to true)
```

### SEE ALSO

* [urlscan scan](urlscan_scan.md)	 - Scan sub-commands

