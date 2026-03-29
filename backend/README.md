
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

## Architecture
The indexer builds an inverse index. 
- It tokenizes the raw messages and maps each token to all the documents it is present in. 
- Each document is identified by it's MsgId.
- Need to create indexes for the categorical fields as well.

The query engine is simply responsible for fetching the documents based on the query.
- Need to implement Union and Intersection algorithms.
- Need to support complex queries that make use of OR and AND operations.
