// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
insert into feeds (id, name, url, created_at, user_id)
values ($1, $2, $3, $4, $5)
returning id, name, url, created_at, user_id, last_fetched_at
`

type CreateFeedParams struct {
	ID        uuid.UUID
	Name      string
	Url       string
	CreatedAt time.Time
	UserID    uuid.UUID
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.Name,
		arg.Url,
		arg.CreatedAt,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.CreatedAt,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeeds = `-- name: GetFeeds :many
select id, name, url, created_at, user_id, last_fetched_at from feeds
`

func (q *Queries) GetFeeds(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.CreatedAt,
			&i.UserID,
			&i.LastFetchedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNextFeedToFetch = `-- name: GetNextFeedToFetch :one
select id, name, url, created_at, user_id, last_fetched_at from feeds
order by last_fetched_at asc nulls first
limit 1
`

func (q *Queries) GetNextFeedToFetch(ctx context.Context) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getNextFeedToFetch)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.CreatedAt,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const updateLastFetchedAt = `-- name: UpdateLastFetchedAt :one
update feeds set last_fetched_at = $1
where id = $2
returning id, name, url, created_at, user_id, last_fetched_at
`

type UpdateLastFetchedAtParams struct {
	LastFetchedAt sql.NullTime
	ID            uuid.UUID
}

func (q *Queries) UpdateLastFetchedAt(ctx context.Context, arg UpdateLastFetchedAtParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, updateLastFetchedAt, arg.LastFetchedAt, arg.ID)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.CreatedAt,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}
