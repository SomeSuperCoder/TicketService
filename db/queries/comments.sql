-- name: CreateComment :one
INSERT INTO comments (
    ticket,
    message
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetCommentsForTicket :many
SELECT
    id,
    ticket,
    message
FROM comments
WHERE ticket = $1
ORDER BY id;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = $1;
