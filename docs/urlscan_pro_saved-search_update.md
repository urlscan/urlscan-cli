## urlscan pro saved-search update

Update a saved search

```
urlscan pro saved-search update [flags]
```

### Examples

```
  urlscan pro save-search update <search-id> -D scans -n <name> -q <query>
  echo "<search-id>" | urlscan pro save-search update - -D scans -n <name> -q <query> -
```

### Options

```
  -D, --datasource string          Which data this Saved Search operates on (hostnames or scans) (required (default "scans")
  -d, --description string         Short description of the saved search (optional)
  -h, --help                       help for update
  -l, --long-description string    Long description of the saved search (optional)
  -n, --name string                Name of the saved search (required)
  -o, --owner-description string   Owner description of the saved search (optional)
  -P, --pass int                   2 for inline-matching, 10 for bookmark-only (required) (default 2)
  -p, --permissions strings        Permissions of the saved search (optional)
  -q, --query string               Search query of the saved search (required)
  -T, --tags strings               Tags of the saved search (optional)
  -t, --tlp string                 TLP (Traffic Light Protocol) of the saved search (required) (default "red")
  -u, --user-tags strings          User tags of the saved search (optional)
```

### SEE ALSO

* [urlscan pro saved-search](urlscan_pro_saved-search.md)	 - Saved search sub-commands

