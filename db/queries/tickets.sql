-- name: CreateTicket :one
INSERT INTO tickets (
    title,
    description,
    location,
    embedding
)
VALUES (
    sqlc.arg(title),
    sqlc.arg(description),
    ST_SetSRID(
        ST_MakePoint(
            sqlc.arg(longitude),
            sqlc.arg(latitude)
        ),
        4326
    )::geography,
    sqlc.arg(embedding)
)
RETURNING *;

-- name: SearchTicketsByMeaning :many
SELECT
    id,
    title,
    description,
    embedding <=> sqlc.arg(query_embedding) AS distance
FROM tickets
ORDER BY distance ASC;
