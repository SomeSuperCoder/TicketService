-- queries/monitoring.sql

-- name: GetOverdueTickets :many
WITH ticket_status_start AS (
    SELECT DISTINCT ON (th.ticket_id)
        th.ticket_id,
        th.created_at as status_start_date
    FROM ticket_history th
    WHERE th.action = 'status_changed'
        AND (th.new_value->>'status' = 'open' OR th.new_value->>'status' = 'init')
    ORDER BY th.ticket_id, th.created_at DESC
),
ticket_created AS (
    SELECT DISTINCT ON (th.ticket_id)
        th.ticket_id,
        th.created_at as created_date
    FROM ticket_history th
    WHERE th.action = 'created'
    ORDER BY th.ticket_id, th.created_at ASC
)
SELECT
    t.id,
    t.description,
    t.status,
    t.subcategory_id,
    t.department_id,
    COALESCE(tss.status_start_date, tc.created_date, t.created_at) as status_start_date,
    GREATEST(0, EXTRACT(DAY FROM CURRENT_DATE - (COALESCE(tss.status_start_date, tc.created_date, t.created_at)::DATE + INTERVAL '7 days')))::INTEGER as lost_days
FROM tickets t
LEFT JOIN ticket_status_start tss ON tss.ticket_id = t.id
LEFT JOIN ticket_created tc ON tc.ticket_id = t.id
WHERE
    t.is_deleted = false
    AND t.is_hidden = false
    AND t.status IN ('open', 'init')
    AND GREATEST(0, EXTRACT(DAY FROM CURRENT_DATE - (COALESCE(tss.status_start_date, tc.created_date, t.created_at)::DATE + INTERVAL '7 days')))::INTEGER >= sqlc.arg('min_lost_days')::INTEGER
    AND (sqlc.narg('department_id')::INTEGER IS NULL OR t.department_id = sqlc.narg('department_id')::INTEGER)
ORDER BY lost_days DESC
LIMIT sqlc.arg('limit')::INTEGER;

-- name: CountOverdueTickets :one
WITH ticket_status_start AS (
    SELECT DISTINCT ON (th.ticket_id)
        th.ticket_id,
        th.created_at as status_start_date
    FROM ticket_history th
    WHERE th.action = 'status_changed'
        AND (th.new_value->>'status' = 'open' OR th.new_value->>'status' = 'init')
    ORDER BY th.ticket_id, th.created_at DESC
),
ticket_created AS (
    SELECT DISTINCT ON (th.ticket_id)
        th.ticket_id,
        th.created_at as created_date
    FROM ticket_history th
    WHERE th.action = 'created'
    ORDER BY th.ticket_id, th.created_at ASC
)
SELECT COUNT(*)
FROM tickets t
LEFT JOIN ticket_status_start tss ON tss.ticket_id = t.id
LEFT JOIN ticket_created tc ON tc.ticket_id = t.id
WHERE
    t.is_deleted = false
    AND t.is_hidden = false
    AND t.status IN ('open', 'init')
    AND GREATEST(0, EXTRACT(DAY FROM CURRENT_DATE - (COALESCE(tss.status_start_date, tc.created_date, t.created_at)::DATE + INTERVAL '7 days')))::INTEGER >= sqlc.arg('min_lost_days')::INTEGER
    AND (sqlc.narg('department_id')::INTEGER IS NULL OR t.department_id = sqlc.narg('department_id')::INTEGER);
