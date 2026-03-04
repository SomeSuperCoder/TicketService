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
  t.id,
  t.status,
  t.description,
  t.is_hidden,
  t.subcategory_id,
  t.department_id,
  t.created_at,
  COALESCE(
    array_agg(tags.id) FILTER (WHERE tags.id IS NOT NULL),
    '{}'
  )::INTEGER[] AS tag_ids,
  COALESCE(
    array_agg(tags.name) FILTER (WHERE tags.name IS NOT NULL),
    '{}'
  )::VARCHAR[] AS tag_names
FROM tickets t
LEFT JOIN ticket_tags ON ticket_tags.ticket = t.id
LEFT JOIN tags ON tags.id = ticket_tags.tag
WHERE t.id = $1 AND t.is_hidden = false AND t.is_deleted = false
GROUP BY t.id;

-- name: GetDetailsForTicket :many
SELECT
  cd.id,
  cd.description,
  cd.sender_name,
  cd.sender_phone,
  cd.sender_email,
  cd.geo_location
FROM complaint_details cd
INNER JOIN tickets t ON t.id = cd.ticket
WHERE cd.ticket = $1
  AND t.is_hidden = false
  AND t.is_deleted = false;

-- name: ListTickets :many
SELECT
  t.id,
  t.status,
  t.description,
  t.is_hidden,
  t.subcategory_id,
  t.department_id,
  t.created_at,
  COALESCE(
    array_agg(tags.id) FILTER (WHERE tags.id IS NOT NULL),
    '{}'
  )::INTEGER[] AS tag_ids,
  COALESCE(
    array_agg(tags.name) FILTER (WHERE tags.name IS NOT NULL),
    '{}'
  )::VARCHAR[] AS tag_names
FROM tickets t
LEFT JOIN ticket_tags ON ticket_tags.ticket = t.id
LEFT JOIN tags ON tags.id = ticket_tags.tag
WHERE
  t.is_hidden = false AND t.is_deleted = false AND
  (sqlc.narg('status')::ticket_status IS NULL OR t.status = sqlc.narg('status')::ticket_status) AND
  (sqlc.narg('subcategory')::INTEGER IS NULL OR t.subcategory_id = sqlc.narg('subcategory')::INTEGER)
GROUP BY t.id
ORDER BY t.embedding <=> sqlc.arg('embedding')::vector
LIMIT sqlc.arg('limit')::INTEGER OFFSET sqlc.arg('offset')::INTEGER;

-- name: UpdateTicket :one
UPDATE tickets
SET 
    status = COALESCE($2, status),
    description = COALESCE($3, description),
    subcategory_id = COALESCE($4, subcategory_id),
    department_id = COALESCE($5, department_id),
    embedding = COALESCE($6, embedding)
WHERE id = $1 AND is_hidden = false AND is_deleted = false
RETURNING *;

-- name: DeleteTicket :one
UPDATE tickets
SET is_deleted = TRUE
WHERE id = $1
RETURNING *;

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
WHERE is_hidden = false AND is_deleted = false
ORDER BY embedding <=> $1
LIMIT $2;

-- name: CountTickets :one
SELECT COUNT(*) FROM tickets
WHERE is_hidden = false AND is_deleted = false;

-- name: CountTicketsByStatus :one
SELECT COUNT(*) FROM tickets
WHERE status = $1 AND is_hidden = false AND is_deleted = false;
