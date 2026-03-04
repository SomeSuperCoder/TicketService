-- name: UpdateDetailsParent :execrows
UPDATE complaint_details
SET ticket = sqlc.arg(original)
WHERE ticket = ANY(sqlc.arg(duplicates)::UUID[]);

-- name: MergeEmbeddings :exec
WITH avg_embedding AS (
    SELECT AVG(embedding) as merged_embedding
    FROM tickets 
    WHERE id = ANY(sqlc.arg(duplcates)::UUID[])
)
UPDATE tickets 
SET embedding = (SELECT merged_embedding FROM avg_embedding),
    status = 'open'
WHERE tickets.id = sqlc.arg(original);

-- name: DeleteAfterMerge :exec
DELETE FROM tickets WHERE id = ANY(sqlc.arg(duplicates)::UUID[]);
