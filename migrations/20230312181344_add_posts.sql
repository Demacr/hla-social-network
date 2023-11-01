-- +goose Up
CREATE TABLE posts (
  id         SERIAL PRIMARY KEY,
  profile_id int NOT NULL,
  title      text,
  text       text NOT NULL,
  FOREIGN KEY (profile_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE posts;
