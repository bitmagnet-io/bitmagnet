-- +goose Up
-- +goose StatementBegin

create table torrent_pieces
(
  info_hash     bytea not null primary key references torrents on delete cascade,
  piece_length  bigint not null,
  pieces        bytea not null,
  created_at    timestamp with time zone not null
);

insert into torrent_pieces (info_hash, piece_length, pieces, created_at)
select info_hash, piece_length, pieces, created_at from torrents where piece_length is not null and pieces is not null;

alter table "torrents" drop column "piece_length";
alter table "torrents" drop column "pieces";

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists torrent_pieces;

alter table "torrents" add column "piece_length" bigint;
alter table "torrents" add column "pieces" bytea;

-- +goose StatementEnd
