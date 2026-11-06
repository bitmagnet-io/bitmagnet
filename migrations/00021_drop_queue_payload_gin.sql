-- +goose NO TRANSACTION
-- +goose Up

-- queue_jobs is polled by queue, status and run_after, and deduplicated by
-- the unique fingerprint index. No query filters or orders by payload, so
-- this GIN index is never used; it only adds write amplification and disk
-- usage. It was introduced with the Postgres queue in 00012.
DROP INDEX CONCURRENTLY IF EXISTS queue_jobs_queue_payload_idx;

-- +goose Down
CREATE INDEX CONCURRENTLY IF NOT EXISTS queue_jobs_queue_payload_idx
    ON queue_jobs USING gin (queue, payload);
