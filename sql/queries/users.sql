-- name: CreateUser :one
insert into users (id, username, api_key, created_at)
values ($1, $2, $3, $4)
returning *;

-- name: GetUserByAPIKey :one
select * from users where api_key = $1;
