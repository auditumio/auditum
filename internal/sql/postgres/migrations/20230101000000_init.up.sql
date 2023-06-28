BEGIN;

CREATE TABLE projects
(
    id                    UUID,
    partition_number      SERIAL      NOT NULL,
    create_time           TIMESTAMPTZ NOT NULL,
    display_name          TEXT        NOT NULL,
    update_record_enabled BOOLEAN,
    delete_record_enabled BOOLEAN,

    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX ON projects (partition_number);

CREATE TABLE records
(
    id                    UUID,
    project_id            UUID        NOT NULL,
    create_time           TIMESTAMPTZ NOT NULL,
    labels                JSONB,
    resource_type         TEXT        NOT NULL,
    resource_id           TEXT        NOT NULL,
    resource_metadata     JSONB,
    operation_type        TEXT        NOT NULL,
    operation_id          TEXT        NOT NULL,
    operation_time        TIMESTAMPTZ NOT NULL,
    operation_metadata    JSONB,
    operation_traceparent TEXT,
    operation_tracestate  TEXT,
    operation_status      SMALLINT,
    actor_type            TEXT        NOT NULL,
    actor_id              TEXT        NOT NULL,
    actor_metadata        JSONB,

    PRIMARY KEY (id, project_id),
    FOREIGN KEY (project_id)
        REFERENCES projects (id)
        ON DELETE RESTRICT
) PARTITION BY LIST (project_id);

CREATE INDEX ON records USING btree (project_id);

-- See: https://www.postgresql.org/docs/current/datatype-json.html#JSON-INDEXING
CREATE INDEX ON records USING gin (labels jsonb_path_ops);

CREATE INDEX ON records USING btree (resource_type);
CREATE INDEX ON records USING btree (resource_id);
CREATE INDEX ON records USING btree (operation_type);
CREATE INDEX ON records USING btree (operation_id);
CREATE INDEX ON records USING btree (operation_time);
CREATE INDEX ON records USING btree (actor_type);
CREATE INDEX ON records USING btree (actor_id);

CREATE TABLE records_resource_changes
(
    record_id   UUID NOT NULL,
    project_id  UUID NOT NULL,
    name        TEXT NOT NULL,
    description TEXT,
    old_value   JSONB,
    new_value   JSONB,

    FOREIGN KEY (record_id, project_id)
        REFERENCES records (id, project_id)
        ON DELETE CASCADE
) PARTITION BY LIST (project_id);

CREATE INDEX ON records_resource_changes USING btree (record_id, project_id);

COMMIT;
