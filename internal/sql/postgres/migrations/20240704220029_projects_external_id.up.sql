BEGIN;

ALTER TABLE projects ADD COLUMN external_id TEXT;

CREATE UNIQUE INDEX idx_projects_external_id ON projects (external_id);

COMMIT;
