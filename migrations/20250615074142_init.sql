-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE recipe_categories
(
    id        VARCHAR(255) PRIMARY KEY,
    name      TEXT NOT NULL UNIQUE,
    image_url TEXT
);

CREATE TABLE ingredients
(
    id   VARCHAR(255) PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    image_url TEXT
);

CREATE TABLE recipes
(
    id             VARCHAR(255) PRIMARY KEY,
    title          TEXT NOT NULL,
    category_id    VARCHAR(255) REFERENCES recipe_categories (id) ON DELETE CASCADE,
    prep_time_min  INT,
    cook_time_min  INT,
    method         TEXT,
    created_at     TIMESTAMP DEFAULT now(),
    ingredient_ids VARCHAR(255)[],
    image_url      TEXT
);

CREATE INDEX idx_recipes_ingredient_ids ON recipes USING GIN (ingredient_ids);

-- +goose Down
DROP INDEX IF EXISTS idx_recipes_ingredient_ids;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS recipe_categories;
