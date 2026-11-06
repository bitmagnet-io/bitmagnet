-- +goose NO TRANSACTION
-- +goose Up

-- +goose StatementBegin
CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS torrent_files_pkey_new
    ON torrent_files (info_hash, "index");
-- +goose StatementEnd

-- +goose StatementBegin
DO $do$
DECLARE
    v_pk_name  name;
    v_pk_cols  text[];
    v_idx_ok   boolean;
    v_idx_cols text[];
    v_con      record;
BEGIN
    -- Serialize concurrent app instances (legacy goose has no session lock).
    -- Do the advisory-lock wait before applying the short AEL timeout.
    PERFORM set_config('lock_timeout', '0', true);
    PERFORM pg_advisory_xact_lock(
        hashtextextended('bitmagnet:migration:torrent_files_key', 0));
    PERFORM set_config('lock_timeout', '5s', true);

    SELECT c.conname, array_agg(a.attname::text ORDER BY k.ord)
      INTO v_pk_name, v_pk_cols
      FROM pg_constraint c
      JOIN LATERAL unnest(c.conkey) WITH ORDINALITY AS k(attnum, ord) ON TRUE
      JOIN pg_attribute a
        ON a.attrelid = c.conrelid AND a.attnum = k.attnum
     WHERE c.conrelid = 'torrent_files'::regclass
       AND c.contype = 'p'
     GROUP BY c.conname;

    IF v_pk_cols IS DISTINCT FROM ARRAY['info_hash', 'index']::text[] THEN
        -- Never attach an invalid / wrong / partial index. Without this the
        -- parser fails with "index ... is not valid" on every retry.
        SELECT i.indisvalid AND i.indisready AND i.indisunique
               AND i.indpred IS NULL AND i.indexprs IS NULL
               AND i.indnkeyatts = 2,
               (SELECT array_agg(a.attname::text ORDER BY k.ord)
                  FROM unnest(string_to_array(i.indkey::text, ' ')::int2[])
                       WITH ORDINALITY AS k(attnum, ord)
                  JOIN pg_attribute a
                    ON a.attrelid = i.indrelid AND a.attnum = k.attnum
                 WHERE k.ord <= i.indnkeyatts)
          INTO v_idx_ok, v_idx_cols
          FROM pg_index i
          JOIN pg_class ic ON ic.oid = i.indexrelid
         WHERE i.indrelid = 'torrent_files'::regclass
           AND ic.relname = 'torrent_files_pkey_new';

        IF NOT FOUND THEN
            RAISE EXCEPTION 'torrent_files_pkey_new is missing; pre-build it first';
        END IF;
        IF NOT v_idx_ok
           OR v_idx_cols IS DISTINCT FROM ARRAY['info_hash', 'index']::text[]
        THEN
            RAISE EXCEPTION 'torrent_files_pkey_new is not a valid unique btree on (info_hash, "index"); run DROP INDEX CONCURRENTLY torrent_files_pkey_new (or REINDEX INDEX CONCURRENTLY) and retry';
        END IF;

        IF v_pk_name IS NOT NULL THEN
            EXECUTE format('ALTER TABLE torrent_files DROP CONSTRAINT %I', v_pk_name);
        END IF;

        -- The constraint takes the pre-built index's name; the conventional
        -- name is applied after the post-condition.
        ALTER TABLE torrent_files
            ADD CONSTRAINT torrent_files_pkey_new
            PRIMARY KEY USING INDEX torrent_files_pkey_new;
    END IF;

    -- Drop any redundant UNIQUE constraint on exactly (info_hash, "index").
    FOR v_con IN
        SELECT c.conname
          FROM pg_constraint c
         WHERE c.conrelid = 'torrent_files'::regclass
           AND c.contype = 'u'
           AND (SELECT array_agg(a.attname::text ORDER BY k.ord)
                  FROM unnest(c.conkey) WITH ORDINALITY AS k(attnum, ord)
                  JOIN pg_attribute a
                    ON a.attrelid = c.conrelid AND a.attnum = k.attnum)
               = ARRAY['info_hash', 'index']::text[]
    LOOP
        EXECUTE format('ALTER TABLE torrent_files DROP CONSTRAINT %I', v_con.conname);
    END LOOP;

    -- Post-condition in the same transaction: bad state rolls back.
    IF (SELECT array_agg(a.attname::text ORDER BY k.ord)
          FROM pg_constraint c
          JOIN LATERAL unnest(c.conkey) WITH ORDINALITY AS k(attnum, ord) ON TRUE
          JOIN pg_attribute a
            ON a.attrelid = c.conrelid AND a.attnum = k.attnum
         WHERE c.conrelid = 'torrent_files'::regclass
           AND c.contype = 'p')
       IS DISTINCT FROM ARRAY['info_hash', 'index']::text[]
    THEN
        RAISE EXCEPTION 'post-condition failed: torrent_files PK is not (info_hash, "index")';
    END IF;

    -- Give the primary key its conventional name. The replacement index is
    -- built under a temporary name because the old primary key owns
    -- torrent_files_pkey until the swap above.
    SELECT c.conname
      INTO v_pk_name
      FROM pg_constraint c
     WHERE c.conrelid = 'torrent_files'::regclass
       AND c.contype = 'p';

    IF v_pk_name IS DISTINCT FROM 'torrent_files_pkey' THEN
        EXECUTE format(
            'ALTER TABLE torrent_files RENAME CONSTRAINT %I TO torrent_files_pkey',
            v_pk_name);
    END IF;
END
$do$;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DO $$
BEGIN
    IF pg_relation_size('torrent_files') > 1024 * 1024 * 1024 THEN
        RAISE EXCEPTION '00022 Down is an offline operation on large tables; restore the previous binary and schema from backup instead';
    END IF;
    IF EXISTS (
        SELECT 1
        FROM torrent_files
        GROUP BY info_hash, path
        HAVING count(*) > 1
    ) THEN
        RAISE EXCEPTION 'cannot restore (info_hash, path) PK: duplicate paths exist';
    END IF;
END;
$$;
-- +goose StatementEnd

-- Build the old indexes. On a small table this is cheap.
CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS torrent_files_pkey_old
    ON torrent_files (info_hash, path);

CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS torrent_files_info_hash_index_key
    ON torrent_files (info_hash, "index");

-- +goose StatementBegin
DO $do$
DECLARE
    v_pk_name name;
    v_pk_cols text[];
BEGIN
    PERFORM set_config('lock_timeout', '5s', true);

    SELECT c.conname, array_agg(a.attname::text ORDER BY k.ord)
      INTO v_pk_name, v_pk_cols
      FROM pg_constraint c
      JOIN LATERAL unnest(c.conkey) WITH ORDINALITY AS k(attnum, ord) ON TRUE
      JOIN pg_attribute a
        ON a.attrelid = c.conrelid AND a.attnum = k.attnum
     WHERE c.conrelid = 'torrent_files'::regclass
       AND c.contype = 'p'
     GROUP BY c.conname;

    IF v_pk_cols IS DISTINCT FROM ARRAY['info_hash', 'path']::text[] THEN
        IF v_pk_name IS NOT NULL THEN
            EXECUTE format('ALTER TABLE torrent_files DROP CONSTRAINT %I', v_pk_name);
        END IF;
        ALTER TABLE torrent_files
            ADD CONSTRAINT torrent_files_pkey
            PRIMARY KEY USING INDEX torrent_files_pkey_old;
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conrelid = 'torrent_files'::regclass
          AND conname = 'torrent_files_info_hash_index_key'
    ) THEN
        ALTER TABLE torrent_files
            ADD CONSTRAINT torrent_files_info_hash_index_key
            UNIQUE USING INDEX torrent_files_info_hash_index_key;
    END IF;
END
$do$;
-- +goose StatementEnd

-- No-op if the old key took the index.
DROP INDEX CONCURRENTLY IF EXISTS torrent_files_pkey_old;
