-- +goose up

create table users (
    id uuid primary key,
    username text unique,
    created_at timestamp not null
);

-- +goose down

drop table users;
