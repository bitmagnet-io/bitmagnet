-- +goose Up
-- +goose StatementBegin

create table key_values
(
  key        text primary key,
  value      text                     not null,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table key_values;

-- +goose StatementEnd
