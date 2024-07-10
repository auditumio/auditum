BEGIN;

ALTER TABLE projects DROP COLUMN external_id;

COMMIT;
