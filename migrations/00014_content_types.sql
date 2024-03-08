-- +goose Up
-- +goose StatementBegin

update torrent_contents set content_type = 'ebook' where content_type = 'book';
update torrent_hints set content_type = 'ebook' where content_type = 'book';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

update torrent_contents set content_type = 'book' where content_type = 'ebook';
update torrent_hints set content_type = 'book' where content_type = 'ebook';

-- +goose StatementEnd
