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

-- name: GetCategoryStatistics :many
SELECT 
    c.id,
    c.name,
    COUNT(t.id) as ticket_count
FROM categories c
LEFT JOIN subcategories sc ON sc.category_id = c.id
LEFT JOIN tickets t ON t.subcategory_id = sc.id 
    AND t.is_deleted = false 
    AND t.is_hidden = false
GROUP BY c.id, c.name
ORDER BY ticket_count DESC, c.name;
