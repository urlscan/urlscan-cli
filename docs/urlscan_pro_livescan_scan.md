## urlscan pro livescan scan

Task a URL to be scanned

```
urlscan pro livescan scan [flags]
```

### Examples

```
  urlscan pro livescan scan <url>
  echo <url> | urlscan pro livescan scan
```

### Options

```
  -b, --blocking                       Whether to do a blocking scan or not (default true)
  -c, --capture-delay int              Delay after page has finished loading before capturing page content (in ms) (default 10000)
  -d, --disable-features strings       Features to disable (annotation, dom, downloads, hideheadless, pageInformation, responses, screenshot, screenshotCompression, stealth)
  -e, --enable-features strings        Features to enable (bannerBypass, downloadWait, fullscreen)
  -H, --extra-headers stringToString   Extra headers to send with the request (e.g., User-Agent: urlscan-cli) (default [])
  -h, --help                           help for scan
  -p, --page-timeout int               Time to wait for the whole scan process (in ms) (default 10000)
  -s, --scanner-id string              ID of the scanner (required)
  -v, --visibility string              Visibility of the scan (public, unlisted or private) (default "private")
```

### SEE ALSO

* [urlscan pro livescan](urlscan_pro_livescan.md)	 - Livescan sub-commands

