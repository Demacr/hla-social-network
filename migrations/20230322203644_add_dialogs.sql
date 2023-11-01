-- +goose Up
CREATE TABLE dialogs (
  id  SERIAL PRIMARY KEY,
  id1 INT NOT NULL,
  id2 INT NOT NULL,
  FOREIGN KEY(id1) REFERENCES users(id),
  FOREIGN KEY(id2) REFERENCES users(id),
  CONSTRAINT id1_less_then_id2 CHECK(id1 < id2)
);

CREATE INDEX dlgs_12_idx on dialogs(id1, id2);
CREATE INDEX dlgs_21_idx on dialogs(id2, id1);

CREATE TABLE messages (
  dialog_id INT NOT NULL,
  id_from   INT NOT NULL,
  id_to     INT NOT NULL,
  seq       INT NOT NULL,
  ts        TIMESTAMP NOT NULL,
  text      VARCHAR(1024) NOT NULL,
  FOREIGN KEY (dialog_id) REFERENCES dialogs (id),
  FOREIGN KEY (id_from)   REFERENCES users   (id),
  FOREIGN KEY (id_to)     REFERENCES users   (id)
);

-- +goose Down
DROP TABLE messages;
DROP INDEX dlgs_21_idx on dialogs;
DROP INDEX dlgs_12_idx on dialogs;
DROP TABLE dialogs;
