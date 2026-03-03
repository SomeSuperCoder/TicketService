-- queries/tickets.sql

-- name: CreateTicket :one
INSERT INTO tickets (
    status,
    complaints,
    description,
    is_hidden,
    subcategory_id,
    department_id,
    embedding
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: CreateTicketWithDefaults :one
INSERT INTO tickets (
    complaints,
    description,
    subcategory_id,
    department_id,
    embedding
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetTicket :one
SELECT * FROM tickets
WHERE id = $1 AND is_hidden = false;

-- name: ListTickets :many
SELECT * FROM tickets
WHERE is_hidden = false
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateTicketStatus :one
UPDATE tickets
SET status = $2
WHERE id = $1 AND is_hidden = false
RETURNING *;

-- name: UpdateTicket :one
UPDATE tickets
SET 
    status = COALESCE($2, status),
    complaints = COALESCE($3, complaints),
    description = COALESCE($4, description),
    subcategory_id = COALESCE($5, subcategory_id),
    department_id = COALESCE($6, department_id),
    embedding = COALESCE($7, embedding)
WHERE id = $1 AND is_hidden = false
RETURNING *;

-- name: HideTicket :exec
UPDATE tickets
SET is_hidden = true
WHERE id = $1;

-- name: DeleteTicket :exec
DELETE FROM tickets
WHERE id = $1;

-- name: SearchTicketsByEmbedding :many
SELECT * FROM tickets
WHERE is_hidden = false
ORDER BY embedding <=> $1
LIMIT $2;

-- name: GetTicketsBySubcategory :many
SELECT * FROM tickets
WHERE subcategory_id = $1 
  AND is_hidden = false
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTicketsByDepartment :many
SELECT * FROM tickets
WHERE department_id = $1 
  AND is_hidden = false
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTicketsByStatus :many
SELECT * FROM tickets
WHERE status = $1 
  AND is_hidden = false
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountTickets :one
SELECT COUNT(*) FROM tickets
WHERE is_hidden = false;

-- name: CountTicketsByStatus :one
SELECT COUNT(*) FROM tickets
WHERE status = $1 AND is_hidden = false;
