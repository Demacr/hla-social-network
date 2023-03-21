-- +goose Up
ALTER TABLE friendship ADD CONSTRAINT id1_not_equal_id2 CHECK(id1 <> id2);
CREATE INDEX id1_id2_idx ON friendship(id1,id2);

-- +goose Down
ALTER TABLE friendship DROP CONSTRAINT id1_not_equal_id2;
DROP INDEX id1_id2_idx;
