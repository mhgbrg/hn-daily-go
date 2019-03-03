ALTER TABLE digest ADD COLUMN start_time timestamptz;
ALTER TABLE digest ADD COLUMN end_time timestamptz;

ALTER TABLE story DROP CONSTRAINT story_external_id_unique;
