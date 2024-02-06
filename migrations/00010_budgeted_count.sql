-- +goose Up
-- +goose StatementBegin

CREATE OR REPLACE FUNCTION budgeted_count(
    query text,
    budget double precision,
    OUT count integer,
    OUT cost double precision,
    OUT budget_exceeded boolean,
    OUT plan jsonb
) LANGUAGE plpgsql AS $$
BEGIN
    EXECUTE 'EXPLAIN (FORMAT JSON) ' || query INTO plan;
    cost := plan->0->'Plan'->'Total Cost';
    IF cost > budget THEN
        count := plan->0->'Plan'->'Plan Rows';
        budget_exceeded := true;
    ELSE
        EXECUTE 'SELECT count(*) FROM (' || query || ') AS subquery' INTO count;
        budget_exceeded := false;
    END IF;
END;
$$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop function if exists budgeted_count(text, double precision, OUT integer, OUT double precision, OUT boolean, OUT jsonb);

-- +goose StatementEnd
