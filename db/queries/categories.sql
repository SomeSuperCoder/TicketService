-- name: GetCategories :many
SELECT * FROM categories;

-- name: GetSubcategoriesFor :many
SELECT * FROM subcategories where category_id = $1;

-- name: CreateCategory :one
INSERT INTO categories (name) VALUES ($1) RETURNING *;

-- name: CreateSubcategory :one
INSERT INTO subcategories (category_id, name) VALUES ($1, $2) RETURNING *;

-- name: CreateDepartment :one
INSERT INTO departments (name) VALUES ($1) RETURNING *;

-- name: GetDepartments :many
SELECT * FROM departments;

-- name: CreateTag :one
INSERT INTO tags (name) VALUES ($1) RETURNING *;

-- name: CreateSource :one
INSERT INTO sources (name) VALUES ($1) RETURNING *;
