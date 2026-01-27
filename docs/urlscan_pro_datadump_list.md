## urlscan pro datadump list

Get the list of data dump files

```
urlscan pro datadump list [flags]
```

### Examples

```
  urlscan pro datadump list days/api
  urlscan pro datadump list hours/api/20260101
  echo "<path>" | urlscan pro datadump list -

  NOTE: path format is <time-window>/<file-type>/<date>
        - time-window: days | hours | minutes. Required.
        - file-type: api | search | screenshot | dom. Required.
        - date: YYYYMMDD format date (optional if time-window is days). Optional.
        if date is not provided, all the available files (files within the last 7 days) will be listed
```

### Options

```
  -h, --help   help for list
```

### SEE ALSO

* [urlscan pro datadump](urlscan_pro_datadump.md)	 - Data dump sub-commands

