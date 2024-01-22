-- +goose Up
-- +goose StatementBegin

-- Fix for pre-0.1.0 users:
delete from content where type='movie' and adult;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
