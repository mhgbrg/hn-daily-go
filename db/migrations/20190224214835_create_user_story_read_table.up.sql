CREATE TABLE user_story_read (
  user_id text NOT NULL,
  story_id int NOT NULL REFERENCES story (id),
  PRIMARY KEY(user_id, story_id)
);
