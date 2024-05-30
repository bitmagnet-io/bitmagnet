-- +goose Up
-- +goose StatementBegin

alter table torrent_contents add column size bigint not null default 0;
alter table torrent_contents add column files_count integer;
create index on torrent_contents(size);
create index on torrent_contents(coalesce(files_count, 0));

create index torrent_contents_seeders_coalesce_idx on torrent_contents (coalesce(seeders, -1));
create index torrent_contents_leechers_coalesce_idx on torrent_contents (coalesce(leechers, -1));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop index torrent_contents_seeders_coalesce_idx;
drop index torrent_contents_leechers_coalesce_idx;

-- +goose StatementEnd
