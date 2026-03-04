-- +goose Up
-- 1. Create independent dictionaries

CREATE TABLE sources (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  CONSTRAINT id_no_zero CHECK (id != 0)
);

CREATE TABLE statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    CONSTRAINT id_no_zero CHECK (id != 0)
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    CONSTRAINT id_no_zero CHECK (id != 0)
);

CREATE TABLE departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    CONSTRAINT id_no_zero CHECK (id != 0)
);

CREATE TABLE tags (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  CONSTRAINT id_no_zero CHECK (id != 0)
);

-- 2. Create dependent dictionaries
CREATE TABLE subcategories (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    CONSTRAINT uq_category_subcategory UNIQUE (category_id, name)
);

-- +goose Down
DROP TABLE subcategories;
DROP TABLE statuses;
DROP TABLE categories;
DROP TABLE departments;
DROP TABLE sources;
DROP TABLE tags;
