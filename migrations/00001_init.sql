-- +goose Up
-- +goose StatementBegin
CREATE
EXTENSION IF NOT EXISTS pg_trgm;

create table torrent_sources
(
  key        text primary key,
  name       text                     not null,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null
);

insert into torrent_sources (key, name, created_at, updated_at)
values ('dht', 'DHT', now(), now()),
       ('rarbg', 'RARBG', now(), now());

create table torrents
(
  info_hash     bytea                    not null primary key,
  name          text                     not null,
  size          bigint                   not null,
  private       boolean                  not null,
  single_file   boolean,
  extension     text generated always as (case
                                            when single_file
                                              then substring(lower(name) from '[^/.]\.([a-z0-9]+)$') end) stored,
  piece_length  bigint,
  pieces        bytea,
  search_string text                     not null,
  tsv           tsvector                 not null generated always as (
    to_tsvector('simple', search_string)
    ) STORED,
  created_at    timestamp with time zone not null,
  updated_at    timestamp with time zone not null
);
create index on torrents (name);
create index on torrents (size);
create index on torrents (single_file);
create index on torrents (extension);
create index on torrents (search_string);
create index on torrents (tsv);
create index on torrents (created_at);

create table torrents_torrent_sources
(
  source       text                     not null references torrent_sources on delete cascade,
  info_hash    bytea                    not null references torrents on delete cascade,
  import_id    text,
  bfsd         bytea,
  bfpe         bytea,
  seeders      integer,
  leechers     integer,
  published_at timestamp with time zone,
  created_at   timestamp with time zone not null,
  updated_at   timestamp with time zone not null,
  primary key (source, info_hash)
);
create index on torrents_torrent_sources (info_hash);
create index on torrents_torrent_sources (import_id);
create index on torrents_torrent_sources (seeders);
create index on torrents_torrent_sources (leechers);
create index on torrents_torrent_sources (created_at);

create table torrent_files
(
  info_hash  bytea                    not null references torrents on delete cascade,
  index      integer                  not null,
  path       text                     not null,
  extension  text generated always as (substring(lower(path) from '[^/.]\.([a-z0-9]+)$')) stored,
  size       bigint                   not null,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null,
  primary key (info_hash, path),
  unique (info_hash, index)
);
create index on torrent_files (size);
create index on torrent_files (extension);

create table metadata_sources
(
  key        text primary key,
  name       text                     not null,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null
);

insert into metadata_sources (key, name, created_at, updated_at)
values ('tmdb', 'TMDB', now(), now()),
       ('imdb', 'IMDB', now(), now()),
       ('tvdb', 'The TVDB', now(), now());

create table content
(
  type              text                     not null,
  source            text                     not null references metadata_sources on delete cascade,
  id                text                     not null,
  title             text                     not null,
  release_date      date,
  release_year      integer,
  adult             boolean,
  original_language text,
  original_title    text,
  overview          text,
  runtime           integer,
  popularity        float,
  vote_average      float,
  vote_count        bigint,
  search_string     text                     not null,
  tsv               tsvector                 not null generated always as (
    to_tsvector('simple', search_string)
    ) stored,
  created_at        timestamp with time zone not null,
  updated_at        timestamp with time zone not null,
  primary key (type, source, id),
  check ((release_date is null) or (release_year = TO_CHAR(release_date, 'yyyy')::INT))
);
create index on content (type);
create index on content (source);
create index on content (id);
create index on content (release_date);
create index on content (adult);
create index on content (original_language);
create index on content (popularity);
CREATE INDEX on content USING gist (search_string gist_trgm_ops);
CREATE INDEX on content USING GIN(tsv);

create table content_attributes
(
  content_type   text                     not null,
  content_source text                     not null references metadata_sources on delete cascade,
  content_id     text                     not null,
  source         text                     not null references metadata_sources on delete cascade,
  key            text                     not null,
  value          text                     not null,
  created_at     timestamp with time zone not null,
  updated_at     timestamp with time zone not null,
  primary key (content_type, content_source, content_id, source, key),
  unique (content_type, content_source, content_id, source, key),
  foreign key (content_type, content_source, content_id) references content (type, source, id) on delete cascade
);
create index on content_attributes (source, key);
create index on content_attributes (source);
create index on content_attributes (key);

create table content_collections
(
  type       text                     not null,
  source     text                     not null references metadata_sources on delete cascade,
  id         text                     not null,
  name       text                     not null,
--   poster_path   text,
--   backdrop_path text,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null,
  primary key (type, source, id)
);
create index on content_collections (type);
create index on content_collections (source);
create index on content_collections (id);
create index on content_collections (name);

create table content_collections_content
(
  content_type              text not null,
  content_source            text not null references metadata_sources (key) on delete cascade,
  content_id                text not null,
  content_collection_type   text not null,
  content_collection_source text not null references metadata_sources on delete cascade,
  content_collection_id     text not null,
  primary key (
               content_type, content_source, content_id, content_collection_type,
               content_collection_source, content_collection_id
    ),
  foreign key (content_type, content_source, content_id) references content (type, source, id) on delete cascade,
  foreign key (
               content_collection_type, content_collection_source,
               content_collection_id
    ) references content_collections (type, source, id) on delete cascade
);
create index on content_collections_content (content_type);
create index on content_collections_content (content_source);
create index on content_collections_content (content_id);
create index on content_collections_content (content_collection_type);
create index on content_collections_content (content_collection_source);
create index on content_collections_content (content_collection_id);

create table torrent_contents
(
  id               text primary key generated always as (encode(info_hash, 'hex') || ':' ||
                                                         coalesce(content_type, '?') ||
                                                         ':' || coalesce(content_source, '?') || ':' ||
                                                         coalesce(content_id, '?')) stored,
  info_hash        bytea                    not null references torrents on delete cascade,
  content_type     text null,
  content_source   text null references metadata_sources on delete cascade,
  content_id       text null,
  title            text                     not null,
  release_date     date,
  release_year     integer,
  external_ids     JSONB,
  languages        JSONB,
  episodes         JSONB,
  video_resolution text,
  video_source     text,
  video_codec      text,
  video_3d         text,
  video_modifier   text,
  release_group    text,
  search_string    text                     not null,
  tsv              tsvector                 not null generated always as (to_tsvector('simple', search_string)) stored,
  created_at       timestamp with time zone not null,
  updated_at       timestamp with time zone not null,
  unique (info_hash, content_type, content_source, content_id),
  foreign key (content_type, content_source, content_id) references content (type, source, id) on delete cascade,
  check ((content_type is not null) or (content_id is null)),
  check ((content_source is null) or (content_id is not null)),
  check ((release_date is null) or (release_year = TO_CHAR(release_date, 'yyyy')::INT))
);
CREATE INDEX on torrent_contents USING gist (id gist_trgm_ops);
create index on torrent_contents (content_type);
create index on torrent_contents (content_source);
create index on torrent_contents (content_id);
create index on torrent_contents (content_source, content_id);
create index on torrent_contents (release_date);
create index on torrent_contents (release_year);
create index on torrent_contents (external_ids);
create index on torrent_contents (languages);
create index on torrent_contents (episodes);
create index on torrent_contents (video_resolution);
create index on torrent_contents (video_source);
create index on torrent_contents (video_codec);
create index on torrent_contents (video_modifier);
create index on torrent_contents (video_3d);
create index on torrent_contents (release_group);
CREATE INDEX on torrent_contents USING gist (search_string gist_trgm_ops);
CREATE INDEX on torrent_contents USING GIN(tsv);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists metadata_sources cascade;
drop table if exists content cascade;
drop table if exists content_attributes cascade;
drop table if exists content_collections cascade;
drop table if exists content_collections_content cascade;
drop table if exists torrent_sources cascade;
drop table if exists torrents cascade;
drop table if exists torrent_files cascade;
drop table if exists torrent_contents cascade;
drop table if exists torrents_torrent_sources cascade;
drop
extension if exists pg_trgm;

-- +goose StatementEnd
