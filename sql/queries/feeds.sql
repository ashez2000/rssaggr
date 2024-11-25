-- name: CreateFeed :one
insert into feeds (id, name, url, created_at, user_id)
values ($1, $2, $3, $4, $5)
returning *;
