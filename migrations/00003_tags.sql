-- +goose Up
-- +goose StatementBegin

create table torrent_tags
(
  info_hash     bytea                    not null references torrents on delete cascade,
  name          text                     not null,
  created_at    timestamp with time zone not null,
  updated_at    timestamp with time zone not null,
  primary key (info_hash, name)
);
create index on torrent_tags (name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists torrent_tags cascade;

-- +goose StatementEnd
