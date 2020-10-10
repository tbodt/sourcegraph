BEGIN;

-- TODO - make an enum instead
ALTER TABLE lsif_nearest_uploads ADD COLUMN forward boolean;
ALTER TABLE lsif_nearest_uploads ADD COLUMN overwritten boolean;
UPDATE lsif_nearest_uploads SET forward = false, overwritten = false;
ALTER TABLE lsif_nearest_uploads ALTER COLUMN forward SET NOT NULL;
ALTER TABLE lsif_nearest_uploads ALTER COLUMN overwritten SET NOT NULL;

-- Mark all repositories as dirty so that we will refresh them
UPDATE lsif_dirty_repositories SET dirty_token = dirty_token + 1;

COMMIT;
