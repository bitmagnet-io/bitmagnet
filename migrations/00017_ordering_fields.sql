-- +goose Up
-- +goose StatementBegin

-- Fix an egregious error that has been there from day 1:
update torrents_torrent_sources set seeders = leechers, leechers = seeders where source = 'dht';

-- Drop some fairly useless columns that take up a lot of space:
alter table torrents_torrent_sources drop column bfsd;
alter table torrents_torrent_sources drop column bfpe;

alter table torrent_contents add column seeders integer;
alter table torrent_contents add column leechers integer;

alter table torrent_contents add column published_at timestamp with time zone;
update torrent_contents set published_at = created_at;
alter table torrent_contents alter column published_at set not null;

create index on torrent_contents (seeders);
create index on torrent_contents (leechers);
create index on torrent_contents (published_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table torrent_contents drop column seeders;
alter table torrent_contents drop column leechers;
alter table torrent_contents drop column published_at;

update torrents_torrent_sources set seeders = leechers, leechers = seeders where source = 'dht';

alter table torrents_torrent_sources add column bfsd bytea;
alter table torrents_torrent_sources add column bfpe bytea;

-- +goose StatementEnd
