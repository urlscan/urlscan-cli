# Analyzing Search Results With DuckDB

A recipe for loading urlscan.io search results in DuckDB.

## Requirements

- [DuckDB](https://duckdb.org/)
- [jq](https://jqlang.org/)

## How To

```bash
urlscan search ... > search.json
# convert search results into JSON Lines
cat search.json  | jq -c '.results[]' > search.jsonl
```

Open Duck DB UI by:

```bash
ducdb -ui
```

Search JSON Lines by:

```sql
SELECT * FROM read_json("search.jsonl")
# search specific field by unnest
SELECT unnest(page) FROM read_json("search.jsonl")
```
