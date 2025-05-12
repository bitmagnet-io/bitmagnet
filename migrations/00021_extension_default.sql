-- +goose Up
-- +goose StatementBegin

alter table "torrents" drop column "extension";
alter table "torrents" add column "extension" text generated always as (case
  when files_status = 'single'::"FilesStatus"
    then substring(lower(name) from '.\.([a-z0-9]+)$') end) stored;
create index on torrents (extension);
alter table "torrent_files" drop column "extension";
alter table "torrent_files" add column "extension" text generated always as 
    (substring(lower(path) from '.\.([a-z0-9]+)$')) stored;
create index on torrent_files (extension);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table "torrents" drop column "extension";
alter table "torrents" add column "extension" text generated always as (case
  when files_status = 'single'::"FilesStatus"
    then substring(lower(name) from '[^/.]\.([a-z0-9]+)$') end) stored;
create index on torrents (extension);
alter table "torrent_files" drop column "extension";
alter table "torrent_files" add column "extension" text generated always as 
    (substring(lower(path) from '[^/.]\.([a-z0-9]+)$')) stored;
create index on torrent_files (extension);

-- +goose StatementEnd
