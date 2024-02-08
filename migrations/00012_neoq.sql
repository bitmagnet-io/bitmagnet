-- +goose Up
-- +goose StatementBegin

CREATE TYPE job_status AS ENUM (
		'new',
		'processed',
		'failed'
);

CREATE SEQUENCE neoq_dead_jobs_id_seq AS bigint;

CREATE TABLE IF NOT EXISTS neoq_dead_jobs (
		id bigint NOT NULL PRIMARY KEY DEFAULT nextval('neoq_dead_jobs_id_seq'::regclass),
		fingerprint text NOT NULL,
		queue text NOT NULL,
		status job_status NOT NULL default 'failed',
		payload jsonb,
		retries integer,
		max_retries integer,
		created_at timestamp with time zone DEFAULT now(),
		error text,
    deadline timestamp with time zone
);

ALTER SEQUENCE neoq_dead_jobs_id_seq OWNED BY neoq_dead_jobs.id;

CREATE SEQUENCE neoq_jobs_id_seq AS bigint;

CREATE TABLE IF NOT EXISTS neoq_jobs (
		id bigint NOT NULL PRIMARY KEY DEFAULT nextval('neoq_jobs_id_seq'::regclass),
		fingerprint text NOT NULL,
		queue text NOT NULL,
		status job_status NOT NULL default 'new',
		payload jsonb,
		retries integer default 0,
		max_retries integer default 23,
		run_after timestamp with time zone DEFAULT now(),
		ran_at timestamp with time zone,
		created_at timestamp with time zone DEFAULT now(),
		error text,
    deadline timestamp with time zone
);

ALTER SEQUENCE neoq_jobs_id_seq OWNED BY neoq_jobs.id;

CREATE INDEX IF NOT EXISTS neoq_job_fetcher_idx ON neoq_jobs (id, status, run_after);
CREATE INDEX IF NOT EXISTS neoq_jobs_fetcher_idx ON neoq_jobs (queue, status, run_after);
CREATE INDEX IF NOT EXISTS neoq_jobs_fingerprint_idx ON neoq_jobs (fingerprint, status);

--- This unique partial index prevents multiple unprocessed jobs with the same payload from being queued
CREATE UNIQUE INDEX IF NOT EXISTS neoq_jobs_fingerprint_unique_idx ON neoq_jobs (fingerprint, status) WHERE NOT (status = 'processed');

CREATE OR REPLACE FUNCTION announce_job() RETURNS trigger AS $$
DECLARE
BEGIN
  PERFORM pg_notify(CAST(NEW.queue AS text), CAST(NEW.id AS text));
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER announce_job
  AFTER INSERT ON neoq_jobs FOR EACH ROW
  WHEN (NEW.run_after <= timezone('utc', NEW.created_at))
  EXECUTE PROCEDURE announce_job();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS neoq_dead_jobs;
DROP TABLE IF EXISTS neoq_jobs;

DROP TYPE IF EXISTS job_status;

-- +goose StatementEnd
