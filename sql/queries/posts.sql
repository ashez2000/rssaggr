-- name: CreatePost :one
insert into posts (id, title, description, url, created_at, published_at, feed_id)
values ($1, $2, $3, $4, $5, $6, $7)
returning *;
