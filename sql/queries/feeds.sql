-- name: CreateFeed :one
insert into feeds (id, name, url, created_at, user_id)
values ($1, $2, $3, $4, $5)
returning *;

-- name: GetFeeds :many
select * from feeds;

-- name: GetNextFeedsToFetch :many
select * from feeds
order by last_fetched_at asc nulls first
limit $1;

-- name: UpdateLastFetchedAt :one
update feeds set last_fetched_at = $1
where id = $2
returning *;
