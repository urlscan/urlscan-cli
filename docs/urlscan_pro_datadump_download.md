## urlscan pro datadump download

Download the data dump file

```
urlscan pro datadump download [flags]
```

### Examples

```
  urlscan pro datadump download days/api/20260101.gz
  urlscan pro datadump download hours/api/20260101/20260101-01.gz
  echo "<path>" | urlscan pro datadump download -

  # use --follow option to download all files from a datadump path
  # for example, the following commands download all the files listed by 'urlscan pro datadump list hours/dom/20260101/'
  # note: --follow memoizes downloaded files in a local database to avoid re-downloading, so it's safe to run it periodically
  urlscan pro datadump download hours/dom/20260101/ --follow
  # if date is not provided, all the available files (files within the last 7 days) will be downloaded
  urlscan pro datadump download hours/api/ --follow
```

### Options

```
  -P, --directory-prefix string   Set directory prefix where file will be saved (default ".")
  -x, --extract                   Extract the downloaded file
  -F, --follow                    Download missing files from the datadump path
  -f, --force                     Force overwrite an existing file
  -h, --help                      help for download
  -o, --output string             Output file name (default <path>.gz)
```

### SEE ALSO

* [urlscan pro datadump](urlscan_pro_datadump.md)	 - Data dump sub-commands

