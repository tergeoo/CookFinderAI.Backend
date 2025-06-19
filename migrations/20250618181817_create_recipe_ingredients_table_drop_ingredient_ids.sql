-- +goose Up
-- +goose StatementBegin
CREATE TABLE recipe_ingredients
(
    recipe_id     VARCHAR(255) REFERENCES recipes (id) ON DELETE CASCADE,
    ingredient_id VARCHAR(255) REFERENCES ingredients (id) ON DELETE CASCADE,
    amount        INTEGER NOT NULL,
    unit          TEXT    NOT NULL, -- 'g', 'ml', 'pcs', etc.
    PRIMARY KEY (recipe_id, ingredient_id)
);

ALTER TABLE recipes
    DROP COLUMN IF EXISTS ingredient_ids;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS recipe_ingredients;
ALTER TABLE recipes
    ADD COLUMN IF NOT EXISTS ingredient_ids TEXT[];
-- +goose StatementEnd
