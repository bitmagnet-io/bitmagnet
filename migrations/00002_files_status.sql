-- +goose Up
-- +goose StatementBegin

create type "FilesStatus" as ENUM ('no_info', 'single', 'multi', 'over_threshold');
alter table "torrents" add column "files_status" "FilesStatus" not null default 'no_info';
update "torrents" set "files_status" = case when single_file then 'single'::"FilesStatus" when single_file = false then 'multi'::"FilesStatus" when single_file is null then 'no_info'::"FilesStatus" end;
create index on torrents (files_status);
alter table "torrents" drop column "extension";
alter table "torrents" drop column "single_file";
alter table "torrents" add column "extension" text generated always as (case
  when files_status = 'single'::"FilesStatus"
    then substring(lower(name) from '[^/.]\.([a-z0-9]+)$') end) stored;
create index on torrents (extension);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table "torrents" add column "single_file" boolean;
update "torrents" set "single_file" = case when "files_status" = 'single'::"FilesStatus" then true when "files_status" = 'multi'::"FilesStatus" then false end;
alter table "torrents" drop column "extension";
alter table "torrents" drop column "files_status";
alter table "torrents" add column "extension" text generated always as (case
  when single_file
    then substring(lower(name) from '[^/.]\.([a-z0-9]+)$') end) stored;
create index on torrents (extension);
drop type "FilesStatus";

-- +goose StatementEnd
