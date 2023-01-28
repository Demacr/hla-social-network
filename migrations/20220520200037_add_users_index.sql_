-- +goose Up
-- +goose StatementBegin
CREATE INDEX fn_ln_idx ON users(name, surname);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX fn_ln_idx;
-- +goose StatementEnd
