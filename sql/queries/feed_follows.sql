-- name: CreateFeedFollow :one
insert into feed_follows (id, created_at, user_id, feed_id)
values ($1, $2, $3, $4)
returning *;

-- name: GetFeedFollows :many
select * from feed_follows where user_id = $1;
