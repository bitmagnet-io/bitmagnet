-- +goose Up
-- +goose StatementBegin

create table content_tags
(
  content_type              text not null,
  content_source            text not null,
  content_id                text not null,
  name          text                     not null,
  created_at    timestamp with time zone not null,
  primary key (content_type, content_source, content_id, name),
  foreign key (content_type, content_source, content_id) references content (type, source, id) on delete cascade,
  check (name ~ '^[a-z0-9]+(-[a-z0-9]+)*$')
);
create index on content_tags (name);
create index on content_tags using gist (name gist_trgm_ops);

alter table torrent_contents add column "tags" text[];
create index on torrent_contents using gin (tags);

alter table torrent_tags drop column updated_at;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


-- +goose StatementEnd
