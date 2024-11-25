-- +goose Up

create table feeds (
    id uuid primary key,
    name text not null,
    url text not null,
    created_at timestamp not null,
    user_id uuid not null references users(id) on delete cascade
);

-- +goose Down

drop table feeds;
