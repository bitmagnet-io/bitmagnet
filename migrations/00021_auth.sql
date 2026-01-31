-- +goose Up
-- +goose StatementBegin

create table roles
(
  name text not null primary key,
  core boolean not null default false,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null
);

insert into roles (name, core, created_at, updated_at) values
  ('admin', true, now(), now()),
  ('editor', true, now(), now()),
  ('user', true, now(), now()),
  ('anon', true, now(), now());

create table role_permissions
(
  role_name text not null references roles(name) on delete cascade,
  namespace text not null,
  object text not null,
  action text not null,
  created_at timestamp with time zone not null,
  primary key (role_name, namespace, object, action)
);

create table users
(
  id serial not null primary key,
  username text not null unique,
  email text unique,
  email_verified boolean not null default false,
  email_verify_code text,
  password bytea,
  role_name text not null references roles(name),
  enabled boolean not null,
  last_login_at timestamp with time zone,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null
);

create table invitations
(
  code text not null primary key,
  role_name text not null default 'user' references roles(name) on delete cascade,
  email text unique,
  created_by integer references users(id) on delete cascade,
  claimed_by integer references users(id) on delete cascade,
  expires_at timestamp with time zone,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null
);

create table api_keys
(
  id serial not null primary key,
  user_id integer references users(id) on delete cascade,
  name text not null,
  hash bytea not null,
  expires_at timestamp with time zone,
  created_at timestamp with time zone not null
);

create table api_key_permissions
(
  api_key_id integer not null references api_keys(id) on delete cascade,
  namespace text not null,
  object text not null,
  action text not null,
  primary key (api_key_id, namespace, object, action)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists api_key_permissions;
drop table if exists api_keys;
drop table if exists invitations;
drop table if exists users;
drop table if exists role_permissions;
drop table if exists permissions;
drop table if exists roles;

-- +goose StatementEnd
