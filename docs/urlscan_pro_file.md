## urlscan pro file

Download a file

```
urlscan pro file [flags]
```

### Examples

```
  urlscan pro file <file-hash>
  echo "<file-hash>" | urlscan pro file -
```

### Options

```
  -F, --filename string   Specify the name of the ZIP file that should be downloaded (defaults to <hash>.zip)
  -f, --force             Enable to force overwriting an existing file
  -h, --help              help for file
  -p, --password string   The password to use to encrypt the ZIP file (default "urlscan!")
```

### SEE ALSO

* [urlscan pro](urlscan_pro.md)	 - Pro sub-commands

