DROP TABLE user_story_read;

CREATE EXTENSION pgcrypto;

CREATE TABLE app_user (
  id serial PRIMARY KEY,
  external_id text NOT NULL UNIQUE,
  first_visit timestamptz NOT NULL
);

CREATE TABLE user_story_read (
  user_id int NOT NULL REFERENCES app_user (id),
  story_id int NOT NULL REFERENCES story (id),
  PRIMARY KEY(user_id, story_id)
);
