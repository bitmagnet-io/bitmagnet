-- +goose Up
-- +goose StatementBegin

alter table queue_jobs add column priority integer not null default 0;

drop index if exists queue_jobs_id_queue_status_run_after_idx;
create index on queue_jobs (id, queue, status, priority, run_after);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table queue_jobs drop column priority;
create index on queue_jobs (id, queue, status, run_after);

-- +goose StatementEnd
