-- name: CreatePost :one
insert into posts (id, title, description, url, created_at, published_at, feed_id)
values ($1, $2, $3, $4, $5, $6, $7)
returning *;

-- name: GetPostsForUser :many
select posts.* from posts
join feed_follows on posts.feed_id = feed_follows.feed_id
where feed_follows.user_id = $1
order by posts.published_at desc
limit $2;
