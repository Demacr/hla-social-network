-- +goose Up
-- +goose StatementBegin
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `surname` varchar(64) NOT NULL,
  `age` int unsigned NOT NULL,
  `sex` varchar(6) NOT NULL,
  `interests` text NOT NULL,
  `city` varchar(64) NOT NULL,
  `email` varchar(128) UNIQUE NOT NULL,
  `password` varchar(64) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `users`;
-- +goose StatementEnd
