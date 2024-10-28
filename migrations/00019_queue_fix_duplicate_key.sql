-- +goose Up
-- +goose StatementBegin

DELETE FROM queue_jobs qj WHERE qj.status = 'retry' AND EXISTS (SELECT * FROM queue_jobs qj2 WHERE qj2.status = 'pending' AND qj2.fingerprint = qj.fingerprint);

DROP INDEX IF EXISTS queue_jobs_fingerprint_status_idx;

CREATE UNIQUE INDEX queue_jobs_fingerprint_idx ON queue_jobs (fingerprint) WHERE status IN ('pending', 'retry');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS queue_jobs_fingerprint_idx;

CREATE UNIQUE INDEX ON queue_jobs (fingerprint, status) WHERE status IN ('pending', 'retry');

-- +goose StatementEnd
