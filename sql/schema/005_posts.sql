-- +goose Up

create table posts (
    id uuid primary key,
    title text not null,
    description text,
    url text not null unique,
    created_at timestamp not null,
    published_at timestamp not null,
    feed_id uuid not null references feeds(id) on delete cascade
);

-- +goose Down

drop table posts;
