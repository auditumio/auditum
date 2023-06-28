BEGIN;

CREATE TABLE projects
(
    id                    UUID,
    partition_number      INTEGER, -- Not used in SQLite, but we keep it to match the common model.
    create_time           TIMESTAMPTZ NOT NULL,
    display_name          TEXT        NOT NULL,
    update_record_enabled BOOLEAN,
    delete_record_enabled BOOLEAN,

    PRIMARY KEY (id)
);

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

    PRIMARY KEY (id),
    FOREIGN KEY (project_id) REFERENCES projects (id)
);

CREATE INDEX idx_records_project_id ON records (project_id);

CREATE INDEX idx_records_labels ON records (labels);

CREATE INDEX idx_records_resource_type ON records (resource_type);
CREATE INDEX idx_records_resource_id ON records (resource_id);

CREATE INDEX idx_records_operation_type ON records (operation_type);
CREATE INDEX idx_records_operation_id ON records (operation_id);
CREATE INDEX idx_records_operation_time ON records (operation_time);

CREATE INDEX idx_records_actor_type ON records (actor_type);
CREATE INDEX idx_records_actor_id ON records (actor_id);

CREATE TABLE records_resource_changes
(
    record_id   UUID NOT NULL,
    project_id  UUID NOT NULL,
    name        TEXT NOT NULL,
    description TEXT,
    old_value   JSONB,
    new_value   JSONB,

    FOREIGN KEY (record_id, project_id) REFERENCES records (id, project_id)
);

CREATE INDEX idx_records_resource_changes_record_id_project_id ON records_resource_changes (record_id, project_id);

COMMIT;
