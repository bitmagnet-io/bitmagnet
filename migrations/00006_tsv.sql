-- +goose Up
-- +goose StatementBegin

alter table "content" drop column "tsv";
alter table "content" drop column "search_string";
alter table "content" add column "tsv" tsvector;
CREATE INDEX on content USING GIN(tsv);

alter table "torrent_contents" drop column "tsv";
alter table "torrent_contents" drop column "search_string";
alter table "torrent_contents" add column "tsv" tsvector;
CREATE INDEX on torrent_contents USING GIN(tsv);

alter table "torrents" drop column "tsv";
alter table "torrents" drop column "search_string";

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table "content" drop column "tsv";
alter table "content" add column "search_string" text not null default '';
alter table "content" add column "tsv" tsvector not null generated always as (to_tsvector('simple', search_string)) stored;
CREATE INDEX on content USING GIN(tsv);

alter table "torrent_contents" drop column "tsv";
alter table "torrent_contents" add column "search_string" text not null default '';
alter table "torrent_contents" add column "tsv" tsvector not null generated always as (to_tsvector('simple', search_string)) stored;
CREATE INDEX on torrent_contents USING GIN(tsv);

alter table "torrents" add column "search_string" text not null default '';
alter table "torrents" add column "tsv" tsvector not null generated always as (to_tsvector('simple', search_string)) stored;
CREATE INDEX on torrents USING GIN(tsv);

-- +goose StatementEnd
