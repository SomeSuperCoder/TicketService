-- name: GetCategories :many
SELECT * FROM categories;

-- name: GetSubcategoriesFor :many
SELECT * FROM subcategories where category_id = $1;
