-- +goose Up
-- +goose StatementBegin

create index if not exists torrent_contents_updated_at_idx on torrent_contents (updated_at);
create index if not exists torrent_contents_content_type_updated_at_idx on torrent_contents (content_type, updated_at);
create index if not exists torrents_torrent_sources_updated_at_idx on torrents_torrent_sources (updated_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop index if exists torrent_contents_updated_at_idx;
drop index if exists torrent_contents_content_type_updated_at_idx;
drop index if exists torrents_torrent_sources_updated_at_idx;

-- +goose StatementEnd
