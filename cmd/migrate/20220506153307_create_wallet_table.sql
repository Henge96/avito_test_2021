-- +goose Up
-- +goose StatementBegin
create type status as enum ('replenishment', 'transfer');

create table wallet
(
    id      serial primary key,
    user_id integer unique               not null,
    balance decimal check (balance >= 0) not null default 0
);

create table transaction
(
    id          serial primary key,
    sender_id   integer   not null,
    receiver_id integer   not null,
    amount      decimal   not null,
    created_at  timestamp not null default now(),
    status      status    not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table transaction;
drop table wallet;
-- +goose StatementEnd
