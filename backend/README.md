
## API
/search 
body : {
    query
}

/upload

## Indexing plan
- Read all parquet files at startup
- Run all messages through the indexer
- Store in memory

```mermaid
flowchart TD
Indexer
QueryEngine
InMemoryStore
```
