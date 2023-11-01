-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id        SERIAL PRIMARY KEY,
  name      varchar(64) NOT NULL,
  surname   varchar(64) NOT NULL,
  age       int NOT NULL,
  sex       varchar(6) NOT NULL,
  interests text NOT NULL,
  city      varchar(64) NOT NULL,
  email     varchar(128) UNIQUE NOT NULL,
  password  varchar(64) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
