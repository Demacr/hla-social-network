-- +goose Up
CREATE TABLE `posts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `profile_id` int NOT NULL,
  `title` text,
  `text` text NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`profile_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE `friendship`;
