-- queries/tickets.sql

-- name: CreateTicketWithDefaults :one
INSERT INTO tickets (
    description,
    complaints,
    subcategory_id,
    department_id,
    embedding
) VALUES (
    sqlc.arg(description),
    ARRAY[ROW(sqlc.arg(sender_name), sqlc.arg(sender_phone), sqlc.arg(sender_email), ST_GeogFromText(sqlc.arg(geo_location)))::complaint_info],
    sqlc.arg(subcategory_id), 
    sqlc.arg(department_id), 
    sqlc.arg(embedding)
) RETURNING *;

-- name: GetTicket :one
SELECT * FROM tickets
WHERE id = $1 AND is_hidden = false;

-- name: ListTickets :many
SELECT * FROM tickets
WHERE is_hidden = false
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

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

-- name: DeleteTicket :exec
DELETE FROM tickets
WHERE id = $1;

-- name: SearchTicketsByEmbedding :many
SELECT * FROM tickets
WHERE is_hidden = false
ORDER BY embedding <=> $1
LIMIT $2;

-- name: CountTickets :one
SELECT COUNT(*) FROM tickets
WHERE is_hidden = false;

-- name: CountTicketsByStatus :one
SELECT COUNT(*) FROM tickets
WHERE status = $1 AND is_hidden = false;
