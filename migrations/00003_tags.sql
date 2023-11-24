-- +goose Up
-- +goose StatementBegin

create table torrent_tags
(
  info_hash     bytea                    not null references torrents on delete cascade,
  name          text                     not null,
  created_at    timestamp with time zone not null,
  updated_at    timestamp with time zone not null,
  primary key (info_hash, name),
  check (name ~ '^[a-z0-9]+(-[a-z0-9]+)*$')
);
create index on torrent_tags (name);
create index on torrent_tags using gist (name gist_trgm_ops);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists torrent_tags cascade;

-- +goose StatementEnd
