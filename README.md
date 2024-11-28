# rssaggr

rssaggr is simple RSS feed aggrefation API where users can create/post their RSS feeds and follow the feeds.
The API aggregate/fetches the RSS feed posts in the background.

> NOTE: project is under development, all features might not work properly

## run locally

> NOTE: run the migration script for setting up database

```bash
cp .env.example .env

# TODO: add build script
go run ./cmd/server/*.go
```

## migration

```bash
cd ./sql/schema

goose postgres postgresql://postgres:password@localhost:5432/rssaggr up
goose postgres postgresql://postgres:password@localhost:5432/rssaggr down
```
