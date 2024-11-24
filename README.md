# rss-aggr

## migration

```bash
cd ./sql/schema

goose postgres postgresql://postgres:password@localhost:5432/rssaggr up
goose postgres postgresql://postgres:password@localhost:5432/rssaggr down
```
