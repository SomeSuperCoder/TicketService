-- +goose Up
-- Create a view that centralizes overdue ticket calculation logic
CREATE OR REPLACE VIEW v_ticket_overdue_status AS
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
    t.id as ticket_id,
    COALESCE(tss.status_start_date, tc.created_date, t.created_at) as status_start_date,
    GREATEST(0, EXTRACT(DAY FROM CURRENT_DATE - (COALESCE(tss.status_start_date, tc.created_date, t.created_at)::DATE + INTERVAL '7 days')))::INTEGER as lost_days,
    CASE 
        WHEN t.status IN ('open', 'init') 
            AND GREATEST(0, EXTRACT(DAY FROM CURRENT_DATE - (COALESCE(tss.status_start_date, tc.created_date, t.created_at)::DATE + INTERVAL '7 days')))::INTEGER > 0
        THEN true
        ELSE false
    END as is_overdue
FROM tickets t
LEFT JOIN ticket_status_start tss ON tss.ticket_id = t.id
LEFT JOIN ticket_created tc ON tc.ticket_id = t.id
WHERE t.is_deleted = false AND t.is_hidden = false;

-- +goose Down
DROP VIEW IF EXISTS v_ticket_overdue_status;
