-- +goose Up
-- +goose StatementBegin

alter table bloom_filters add column oid oid;
update bloom_filters set oid = lo_from_bytea(0, bytes::bytea);
alter table bloom_filters drop column bytes;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table bloom_filters add column bytes bytea;
update bloom_filters set bytes = ''::bytea where true;
alter table bloom_filters alter column bytes set not null;
alter table bloom_filters drop column oid;

-- +goose StatementEnd
