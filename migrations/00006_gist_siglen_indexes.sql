-- +goose Up
-- +goose StatementBegin

create index if not exists torrent_contents_search_string_344_idx on torrent_contents using gist (search_string gist_trgm_ops(siglen=344));
create index if not exists content_search_string_168_idx on content using gist (search_string gist_trgm_ops(siglen=168));

drop index if exists torrent_contents_search_string_idx;
drop index if exists content_search_string_idx;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop index if exists torrent_contents_search_string_344_idx;
drop index if exists content_search_string_168_idx;

create index if not exists torrent_contents_search_string_idx on torrent_contents using gist (id gist_trgm_ops);
create index if not exists content_search_string_idx on content using gist (id gist_trgm_ops);

-- +goose StatementEnd
