-- +goose Up
-- +goose StatementBegin
CREATE TABLE `auth` (
  `user_id` INT NOT NULL,
  `token` VARCHAR(64) NOT NULL,
  `timestamp` TIMESTAMP NOT NULL,
  FOREIGN KEY (`user_id`)
    REFERENCES `users`(`id`)
    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `auth`;
-- +goose StatementEnd
