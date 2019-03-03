ALTER TABLE digest DROP COLUMN start_time;
ALTER TABLE digest DROP COLUMN end_time;

ALTER TABLE story ADD CONSTRAINT story_external_id_unique UNIQUE (external_id);
