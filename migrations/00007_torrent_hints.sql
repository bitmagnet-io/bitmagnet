-- +goose Up
-- +goose StatementBegin

create table torrent_hints
(
  info_hash        bytea not null primary key references torrents on delete cascade,
  content_type     text not null,
  content_source   text null references metadata_sources on delete cascade,
  content_id       text null,
  title            text null,
  release_year     integer null,
  languages        JSONB null,
  episodes         JSONB null,
  video_resolution text null,
  video_source     text null,
  video_codec      text null,
  video_3d         text null,
  video_modifier   text null,
  release_group    text null,
  created_at       timestamp with time zone not null,
  updated_at       timestamp with time zone not null,
  check ((content_source is null) or (content_id is not null))
);

alter table "torrent_contents" drop column "title";
alter table "torrent_contents" drop column "external_ids";
alter table "torrent_contents" drop column "release_date";
alter table "torrent_contents" drop column "release_year";

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table "torrent_contents" add column "title" text not null;
alter table "torrent_contents" add column "external_ids" JSONB;
alter table "torrent_contents" add column "release_date" timestamp with time zone;
alter table "torrent_contents" add column "release_year" integer;

drop table torrent_hints;

-- +goose StatementEnd
