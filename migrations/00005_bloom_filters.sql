-- +goose Up
-- +goose StatementBegin

create table bloom_filters
(
  key           text                  primary key,
  bytes         bytea                    not null,
  created_at    timestamp with time zone not null,
  updated_at    timestamp with time zone not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists bloom_filters;

-- +goose StatementEnd
