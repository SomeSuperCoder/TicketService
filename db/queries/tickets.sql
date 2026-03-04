-- queries/tickets.sql

-- name: CreateTicketWithDefaults :one
INSERT INTO tickets (
    description,
    subcategory_id,
    department_id,
    embedding
) VALUES (
    sqlc.arg(description),
    sqlc.arg(subcategory_id), 
    sqlc.arg(department_id), 
    sqlc.arg(embedding)
) RETURNING *;

-- name: CreateComplaint :one
INSERT INTO complaint_details (
  ticket,
  description,
  sender_name,
  sender_phone,
  sender_email,
  geo_location
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetTicket :one
SELECT
  id,
  status,
  description,
  is_hidden,
  subcategory_id,
  department_id,
  created_at
FROM tickets
WHERE id = $1 AND is_hidden = false;

-- name: GetDetailsForTicket :many
SELECT
  id,
  description,
  sender_name,
  sender_phone,
  sender_email,
  geo_location
FROM complaint_details
WHERE ticket = $1;
  

-- name: ListTickets :many
SELECT * FROM tickets
WHERE is_hidden = false
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateTicket :one
UPDATE tickets
SET 
    status = COALESCE($2, status),
    description = COALESCE($3, description),
    subcategory_id = COALESCE($4, subcategory_id),
    department_id = COALESCE($5, department_id),
    embedding = COALESCE($6, embedding)
WHERE id = $1 AND is_hidden = false
RETURNING *;

-- name: DeleteTicket :exec
DELETE FROM tickets
WHERE id = $1;

-- name: SearchTicketsByEmbedding :many
SELECT 
  id,
  status,
  description,
  is_hidden,
  subcategory_id,
  department_id,
  created_at
FROM tickets
WHERE is_hidden = false
ORDER BY embedding <=> $1
LIMIT $2;

-- name: CountTickets :one
SELECT COUNT(*) FROM tickets
WHERE is_hidden = false;

-- name: CountTicketsByStatus :one
SELECT COUNT(*) FROM tickets
WHERE status = $1 AND is_hidden = false;
