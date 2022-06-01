-- +goose Up
-- +goose StatementBegin
create table wallet
(
    id serial primary key,
    user_id integer unique not null,
    balance decimal check (balance >= 0.0)
);
create table transaction
(
    id serial primary key,
    wallet_id integer references wallet (id),
    receiver_wallet_id integer not null,
    money decimal check(money >= 0.0),
    date timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table transaction;
drop table wallet;
-- +goose StatementEnd
