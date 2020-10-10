BEGIN;

ALTER TABLE lsif_nearest_uploads DROP COLUMN forward;
ALTER TABLE lsif_nearest_uploads DROP COLUMN overwritten;

COMMIT;
