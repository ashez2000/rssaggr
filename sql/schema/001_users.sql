-- +goose Up

create table users (
    id uuid primary key,
    username text unique not null,
    api_key text unique not null,
    created_at timestamp not null
);

-- +goose Down

drop table users;
