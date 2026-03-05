-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    role,
    status,
    department_id,
    first_name,
    last_name,
    middle_name
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 AND status != 'blocked';

-- name: ListUsers :many
SELECT 
    u.*,
    d.name as department_name
FROM users u
LEFT JOIN departments d ON u.department_id = d.id
WHERE 
    (sqlc.narg('role')::user_role IS NULL OR u.role = sqlc.narg('role')::user_role)
    AND (sqlc.narg('status')::user_status IS NULL OR u.status = sqlc.narg('status')::user_status)
    AND (sqlc.narg('email')::TEXT IS NULL OR u.email ILIKE '%' || sqlc.narg('email')::TEXT || '%')
ORDER BY u.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users
WHERE 
    (sqlc.narg('role')::user_role IS NULL OR role = sqlc.narg('role')::user_role)
    AND (sqlc.narg('status')::user_status IS NULL OR status = sqlc.narg('status')::user_status)
    AND (sqlc.narg('email')::TEXT IS NULL OR email ILIKE '%' || sqlc.narg('email')::TEXT || '%');

-- name: UpdateUser :one
UPDATE users
SET
    role = COALESCE(sqlc.narg('role')::user_role, role),
    status = COALESCE(sqlc.narg('status')::user_status, status),
    department_id = COALESCE(sqlc.narg('department_id')::INT, department_id),
    first_name = COALESCE(sqlc.narg('first_name')::VARCHAR, first_name),
    last_name = COALESCE(sqlc.narg('last_name')::VARCHAR, last_name),
    middle_name = COALESCE(sqlc.narg('middle_name')::VARCHAR, middle_name),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING *;
