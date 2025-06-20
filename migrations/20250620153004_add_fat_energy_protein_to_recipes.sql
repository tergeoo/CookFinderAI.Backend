-- +goose Up
-- +goose StatementBegin
ALTER TABLE recipes
    ADD COLUMN fat     FLOAT NOT NULL DEFAULT 0,
    ADD COLUMN energy  INT   NOT NULL DEFAULT 0,
    ADD COLUMN protein FLOAT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE recipes
    DROP COLUMN fat,
    DROP COLUMN energy,
    DROP COLUMN protein;
-- +goose StatementEnd
