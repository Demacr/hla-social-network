-- +goose Up
CREATE TABLE `dialogs` (
  `id`  INT NOT NULL AUTO_INCREMENT,
  `id1` INT NOT NULL,
  `id2` INT NOT NULL,
  PRIMARY KEY(`id`),
  FOREIGN KEY(`id1`) REFERENCES users(`id`),
  FOREIGN KEY(`id2`) REFERENCES users(`id`),
  INDEX(`id1`),
  INDEX(`id2`),
  CONSTRAINT id1_less_then_id2 CHECK(id1 < id2)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `messages` (
  `dialog_id` INT NOT NULL,
  `id_from`   INT NOT NULL,
  `id_to`     INT NOT NULL,
  `seq`       INT NOT NULL,
  `ts`        TIMESTAMP NOT NULL,
  `text`      VARCHAR(1024) NOT NULL,
  FOREIGN KEY (`dialog_id`) REFERENCES `dialogs` (`id`),
  FOREIGN KEY (`id_from`)   REFERENCES `users`   (`id`),
  FOREIGN KEY (`id_to`)     REFERENCES `users`   (`id`),
  PRIMARY KEY (`dialog_id`, `seq`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE `messages`;
DROP TABLE `dialogs`;
