-- +goose Up
-- +goose StatementBegin

drop index if exists torrent_contents_id_idx;
drop index if exists torrent_contents_languages_idx;
drop index if exists torrent_contents_tsv_idx;

create extension if not exists btree_gin;

create index on torrent_contents using gin(content_type, tsv);
create index on torrent_contents using gin(content_type, languages);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop index if exists torrent_contents_content_type_tsv_idx;
drop index if exists torrent_contents_content_type_languages_idx;

CREATE INDEX on torrent_contents USING gist (id gist_trgm_ops);
create index on torrent_contents (languages);
CREATE INDEX on torrent_contents USING GIN(tsv);

drop extension if exists btree_gin;

-- +goose StatementEnd
