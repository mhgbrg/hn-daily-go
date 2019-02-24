CREATE TABLE digest (
  id serial PRIMARY KEY,
  date date NOT NULL,
  start_time timestamptz NOT NULL,
  end_time timestamptz NOT NULL,
  generated_at timestamptz NOT NULL
);

CREATE TABLE story (
  id serial PRIMARY KEY,
  external_id int NOT NULL,
  posted_at timestamptz NOT NULL,
  title text NOT NULL,
  url text NOT NULL,
  author text NOT NULL,
  points text NOT NULL,
  num_comments int NOT NULL,
  digest_id int NOT NULL REFERENCES digest (id)
);
