-- name: CreateFeed :one
insert into feeds (id, name, url, created_at)
values ($1, $2, $3, $4)
returning *;
