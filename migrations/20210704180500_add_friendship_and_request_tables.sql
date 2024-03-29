-- +goose Up
CREATE TABLE friendship (
  id1 int NOT NULL,
  id2 int NOT NULL,
  FOREIGN KEY (id1) REFERENCES users (id),
  FOREIGN KEY (id2) REFERENCES users (id)
);

CREATE TABLE friendship_requests (
  id_from int NOT NULL,
  id_to   int NOT NULL,
  FOREIGN KEY (id_from) REFERENCES users (id),
  FOREIGN KEY (id_to)   REFERENCES users (id)
);

-- +goose Down
DROP TABLE friendship;
DROP TABLE friendship_requests;
