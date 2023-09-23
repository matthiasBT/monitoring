-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
    CREATE TYPE metric_type AS ENUM('gauge', 'counter');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
CREATE TABLE metrics (id text primary key, mtype metric_type, delta bigint, val double precision);
CREATE INDEX IF NOT EXISTS search_idx ON metrics USING btree(id, mtype);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX search_idx;
DROP TABLE metrics;
DROP TYPE metric_type;
-- +goose StatementEnd
