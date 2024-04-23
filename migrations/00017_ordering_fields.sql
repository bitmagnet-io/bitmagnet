-- +goose Up
-- +goose StatementBegin

-- Drop some fairly useless columns that take up a lot of space:
alter table torrents_torrent_sources drop column bfsd, drop column bfpe;

alter table torrent_contents add column seeders integer;
alter table torrent_contents add column leechers integer;

alter table torrent_contents add column published_at timestamp with time zone not null default '1999-01-01 00:00:00+00';

create index on torrent_contents (seeders);
create index on torrent_contents (leechers);
create index on torrent_contents (published_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table torrent_contents drop column seeders;
alter table torrent_contents drop column leechers;
alter table torrent_contents drop column published_at;

alter table torrents_torrent_sources add column bfsd bytea;
alter table torrents_torrent_sources add column bfpe bytea;

-- +goose StatementEnd
