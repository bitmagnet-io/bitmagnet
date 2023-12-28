-- +goose Up
-- +goose StatementBegin

CREATE OR REPLACE FUNCTION tsv_parts_to_tsv(parts JSONB) returns tsvector as $$
select setweight(to_tsvector(coalesce(parts->>'A', '')), 'A') ||
       setweight(to_tsvector(coalesce(parts->>'B', '')), 'B') ||
       setweight(to_tsvector(coalesce(parts->>'C', '')), 'C') ||
       setweight(to_tsvector(coalesce(parts->>'D', '')), 'D')
         $$ language sql immutable;

alter table "torrent_contents" drop column "tsv";
alter table "torrent_contents" drop column "search_string";
alter table "torrent_contents" add column "tsv_parts" JSONB;
alter table "torrent_contents" add column "tsv" tsvector not null generated always as (tsv_parts_to_tsv(tsv_parts)) stored;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table "torrent_contents" drop column "tsv";
alter table "torrent_contents" drop column "tsv_parts";
alter table "torrent_contents" add column "search_string" text not null default '';
alter table "torrent_contents" add column "tsv" tsvector not null generated always as (to_tsvector('simple', search_string)) stored;

drop function tsv_parts_to_tsv(JSONB);

-- +goose StatementEnd
