-- name: CreateUser :one
insert into users (id, username, created_at)
values ($1, $2, $3)
returning *;
