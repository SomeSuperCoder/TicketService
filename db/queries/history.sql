-- name: CreateHistoryEntry :one
INSERT INTO ticket_history (
    ticket_id,
    action,
    old_value,
    new_value,
    user_id,
    user_name,
    user_email,
    description
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetTicketHistory :many
SELECT
    id,
    ticket_id,
    action,
    old_value,
    new_value,
    user_id,
    user_name,
    user_email,
    description,
    created_at
FROM ticket_history
WHERE ticket_id = $1
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')::INTEGER
OFFSET sqlc.arg('offset')::INTEGER;

-- name: GetTicketHistoryCount :one
SELECT COUNT(*) FROM ticket_history WHERE ticket_id = $1;

-- name: GetRecentHistory :many
SELECT
    h.id,
    h.ticket_id,
    h.action,
    h.old_value,
    h.new_value,
    h.user_id,
    h.user_name,
    h.user_email,
    h.description,
    h.created_at,
    t.description as ticket_description
FROM ticket_history h
INNER JOIN tickets t ON t.id = h.ticket_id
WHERE t.is_deleted = false
ORDER BY h.created_at DESC
LIMIT sqlc.arg('limit')::INTEGER
OFFSET sqlc.arg('offset')::INTEGER;

-- name: CreateHistoryEntryWithTime :one
INSERT INTO ticket_history (
    ticket_id,
    action,
    old_value,
    new_value,
    user_name,
    user_email,
    description,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;
