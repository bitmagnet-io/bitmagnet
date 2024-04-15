-- +goose Up
-- +goose StatementBegin

alter table torrent_contents add column seeders integer;
alter table torrent_contents add column leechers integer;
alter table torrent_contents add column published_at timestamp with time zone;
create index on torrent_contents (seeders);
create index on torrent_contents (leechers);
create index on torrent_contents (published_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table torrent_contents drop column seeders;
alter table torrent_contents drop column leechers;
alter table torrent_contents drop column published_at;

-- +goose StatementEnd
