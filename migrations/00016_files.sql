-- +goose Up
-- +goose StatementBegin

alter table torrents add column files_count integer;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table torrents drop column files_count;

-- +goose StatementEnd
