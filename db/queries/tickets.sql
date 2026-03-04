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

-- name: UpdateTicketSimple :one
UPDATE tickets
SET
  status = coalesce(sqlc.narg('status')::ticket_status, status),
  department_id = coalesce(sqlc.narg('department_id')::INTEGER, department_id)
WHERE is_hidden = false AND is_deleted = false AND
  id = sqlc.arg(id)
RETURNING status, department_id;

-- name: DeleteTagsFromTicket :execrows
DELETE FROM ticket_tags
WHERE ticket = sqlc.arg(ticket) AND tag = ANY(sqlc.arg(tags)::INTEGER[]);

-- name: AddTagsToTicket :execrows
INSERT INTO ticket_tags (ticket, tag)
SELECT
  sqlc.arg(ticket),
  unnest(sqlc.arg(tags)::INTEGER[])
ON CONFLICT DO NOTHING;

-- name: DeleteTicket :one
UPDATE tickets
SET is_deleted = true
WHERE id = $1
RETURNING *;

-- name: CountTickets :one
SELECT COUNT(*) FROM tickets WHERE is_hidden = false AND is_deleted = false;
