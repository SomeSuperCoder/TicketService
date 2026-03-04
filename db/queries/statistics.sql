-- name: GetStatisticsSummary :one
SELECT
    COUNT(*) FILTER (WHERE is_deleted = false AND is_hidden = false) AS total,

    COUNT(*) FILTER (
        WHERE status = 'closed'
        AND is_deleted = false
        AND is_hidden = false
    ) AS resolved,

    COUNT(*) FILTER (
        WHERE status = 'open'
        AND is_deleted = false
        AND is_hidden = false
    ) AS in_progress,

    COUNT(DISTINCT department_id) FILTER (
        WHERE department_id IS NOT NULL
        AND is_deleted = false
        AND is_hidden = false
    ) AS active_executors
FROM tickets;
