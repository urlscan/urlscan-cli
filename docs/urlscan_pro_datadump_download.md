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
```

### Options

```
  -x, --extract         Extract the downloaded file
  -f, --force           Force overwrite an existing file
  -h, --help            help for download
  -o, --output string   Output file name (default <path>.gz)
```

### SEE ALSO

* [urlscan pro datadump](urlscan_pro_datadump.md)	 - Data dump sub-commands

