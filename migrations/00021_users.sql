-- +goose Up
-- +goose StatementBegin

create table users
(
  id serial not null primary key,
  username text not null unique,
  password text,
  last_login_at timestamp with time zone,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists users;

-- +goose StatementEnd
