-- +goose Up
-- +goose StatementBegin

CREATE TYPE queue_job_status AS ENUM (
  'pending',
  'processed',
  'retry',
  'failed'
);

CREATE TABLE IF NOT EXISTS queue_jobs (
  id text NOT NULL PRIMARY KEY default gen_random_uuid(),
  fingerprint text NOT NULL,
  queue text NOT NULL,
  status queue_job_status NOT NULL default 'pending',
  payload jsonb NOT NULL,
  retries integer NOT NULL default 0,
  max_retries integer NOT NULL default 0,
  run_after timestamp with time zone not null,
  ran_at timestamp with time zone,
  error text,
  deadline timestamp with time zone,
  archival_duration interval not null,
  created_at timestamp with time zone not null
);

CREATE INDEX ON queue_jobs (queue, status);
CREATE INDEX ON queue_jobs (id, queue, status, run_after);
CREATE INDEX ON queue_jobs USING gin(queue, payload);

--- This unique partial index prevents multiple unprocessed jobs with the same payload from being queued
CREATE UNIQUE INDEX ON queue_jobs (fingerprint, status) WHERE NOT (status IN ('processed', 'failed'));

CREATE OR REPLACE FUNCTION queue_announce_job() RETURNS trigger AS $$
DECLARE
BEGIN
  PERFORM pg_notify(CAST(NEW.queue AS text), CAST(NEW.id AS text));
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER queue_announce_job
  AFTER INSERT ON queue_jobs FOR EACH ROW
  WHEN (NEW.run_after <= NEW.created_at)
  EXECUTE PROCEDURE queue_announce_job();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS queue_jobs;
DROP TYPE IF EXISTS queue_job_status;

-- +goose StatementEnd
