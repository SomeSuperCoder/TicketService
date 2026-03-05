-- queries/monitoring.sql

-- name: GetOverdueTickets :many
SELECT
    t.id,
    t.description,
    t.status,
    t.subcategory_id,
    t.department_id,
    v.status_start_date,
    v.lost_days
FROM tickets t
INNER JOIN v_ticket_overdue_status v ON v.ticket_id = t.id
WHERE
    t.status IN ('open', 'init')
    AND v.lost_days >= sqlc.arg('min_lost_days')::INTEGER
    AND (sqlc.narg('department_id')::INTEGER IS NULL OR t.department_id = sqlc.narg('department_id')::INTEGER)
ORDER BY v.lost_days DESC
LIMIT sqlc.arg('limit')::INTEGER;

-- name: CountOverdueTickets :one
SELECT COUNT(*)
FROM tickets t
INNER JOIN v_ticket_overdue_status v ON v.ticket_id = t.id
WHERE
    t.status IN ('open', 'init')
    AND v.lost_days >= sqlc.arg('min_lost_days')::INTEGER
    AND (sqlc.narg('department_id')::INTEGER IS NULL OR t.department_id = sqlc.narg('department_id')::INTEGER);

-- name: GetDepartmentEfficiency :many
WITH ticket_closed AS (
    SELECT DISTINCT ON (th.ticket_id)
        th.ticket_id,
        th.created_at as closed_date
    FROM ticket_history th
    WHERE th.action = 'status_changed'
        AND th.new_value->>'status' = 'closed'
    ORDER BY th.ticket_id, th.created_at DESC
),
ticket_created AS (
    SELECT DISTINCT ON (th.ticket_id)
        th.ticket_id,
        th.created_at as created_date
    FROM ticket_history th
    WHERE th.action = 'created'
    ORDER BY th.ticket_id, th.created_at ASC
),
department_stats AS (
    SELECT
        d.id as department_id,
        d.name as department_name,
        COUNT(CASE WHEN t.status IN ('open', 'init') THEN 1 END) as in_progress,
        COUNT(CASE WHEN v.is_overdue THEN 1 END) as overdue,
        AVG(CASE 
            WHEN t.status = 'closed' AND tcl.closed_date IS NOT NULL
            THEN EXTRACT(EPOCH FROM (tcl.closed_date - COALESCE(tc.created_date, t.created_at))) / 86400.0
        END) as avg_resolution_days
    FROM departments d
    LEFT JOIN tickets t ON t.department_id = d.id AND t.is_deleted = false AND t.is_hidden = false
    LEFT JOIN v_ticket_overdue_status v ON v.ticket_id = t.id
    LEFT JOIN ticket_created tc ON tc.ticket_id = t.id
    LEFT JOIN ticket_closed tcl ON tcl.ticket_id = t.id
    GROUP BY d.id, d.name
),
department_trend AS (
    SELECT
        d.id as department_id,
        COUNT(CASE 
            WHEN t.status = 'closed' 
            AND tcl.closed_date >= CURRENT_DATE - INTERVAL '7 days'
            THEN 1 
        END) as closed_last_week,
        COUNT(CASE 
            WHEN t.status = 'closed' 
            AND tcl.closed_date >= CURRENT_DATE - INTERVAL '14 days'
            AND tcl.closed_date < CURRENT_DATE - INTERVAL '7 days'
            THEN 1 
        END) as closed_prev_week
    FROM departments d
    LEFT JOIN tickets t ON t.department_id = d.id AND t.is_deleted = false AND t.is_hidden = false
    LEFT JOIN ticket_closed tcl ON tcl.ticket_id = t.id
    GROUP BY d.id
)
SELECT
    ds.department_id,
    ds.department_name,
    ds.in_progress,
    ds.overdue,
    COALESCE(ds.avg_resolution_days, 0)::FLOAT as avg_resolution_days,
    CASE
        WHEN dt.closed_prev_week = 0 THEN 0
        ELSE ((dt.closed_last_week::FLOAT - dt.closed_prev_week::FLOAT) / dt.closed_prev_week::FLOAT * 100)
    END as trend_percent
FROM department_stats ds
LEFT JOIN department_trend dt ON dt.department_id = ds.department_id
ORDER BY ds.department_name;

-- name: GetKPI :one
WITH period_filter AS (
    SELECT 
        CASE 
            WHEN sqlc.arg('period')::TEXT = 'week' THEN CURRENT_DATE - INTERVAL '7 days'
            WHEN sqlc.arg('period')::TEXT = 'month' THEN CURRENT_DATE - INTERVAL '30 days'
            WHEN sqlc.arg('period')::TEXT = 'year' THEN CURRENT_DATE - INTERVAL '365 days'
            ELSE CURRENT_DATE - INTERVAL '30 days'
        END as start_date
),
ticket_created AS (
    SELECT DISTINCT ON (th.ticket_id)
        th.ticket_id,
        th.created_at as created_date
    FROM ticket_history th
    WHERE th.action = 'created'
    ORDER BY th.ticket_id, th.created_at ASC
),
first_response AS (
    SELECT DISTINCT ON (th.ticket_id)
        th.ticket_id,
        th.created_at as first_response_date
    FROM ticket_history th
    WHERE th.action = 'comment_added'
    ORDER BY th.ticket_id, th.created_at ASC
),
avg_response AS (
    SELECT 
        COALESCE(AVG(EXTRACT(EPOCH FROM (fr.first_response_date - tc.created_date)) / 86400.0), 0) as avg_days
    FROM tickets t
    INNER JOIN ticket_created tc ON tc.ticket_id = t.id
    CROSS JOIN period_filter pf
    INNER JOIN first_response fr ON fr.ticket_id = t.id
    WHERE 
        t.is_deleted = false 
        AND t.is_hidden = false
        AND tc.created_date >= pf.start_date
),
overdue_count AS (
    SELECT COUNT(*) as total
    FROM tickets t
    INNER JOIN v_ticket_overdue_status v ON v.ticket_id = t.id
    WHERE t.status IN ('open', 'init') AND v.is_overdue = true
),
satisfaction_data AS (
    SELECT 
        COUNT(CASE WHEN t.status = 'closed' THEN 1 END) as closed_count,
        COUNT(CASE WHEN t.status IN ('open', 'init') AND v.is_overdue = false THEN 1 END) as on_time_count
    FROM tickets t
    INNER JOIN ticket_created tc ON tc.ticket_id = t.id
    CROSS JOIN period_filter pf
    LEFT JOIN v_ticket_overdue_status v ON v.ticket_id = t.id
    WHERE 
        t.is_deleted = false 
        AND t.is_hidden = false
        AND tc.created_date >= pf.start_date
)
SELECT
    ar.avg_days::FLOAT as avg_response_days,
    oc.total as overdue_count,
    CASE 
        WHEN (sd.closed_count + sd.on_time_count) = 0 THEN 0
        ELSE ((sd.closed_count::FLOAT + sd.on_time_count::FLOAT * 0.5) / (sd.closed_count + sd.on_time_count) * 100)
    END as satisfaction_index
FROM avg_response ar
CROSS JOIN overdue_count oc
CROSS JOIN satisfaction_data sd;
