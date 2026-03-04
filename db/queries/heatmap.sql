-- name: GetHeatmapPoints :many
SELECT
    ST_Y(cd.geo_location::geometry)::DOUBLE PRECISION as lat,
    ST_X(cd.geo_location::geometry)::DOUBLE PRECISION as lng,
    COUNT(*) as count,
    COUNT(*) as intensity
FROM complaint_details cd
INNER JOIN tickets t ON t.id = cd.ticket
WHERE 
    t.is_deleted = false 
    AND t.is_hidden = false
    AND t.created_at >= sqlc.arg('start_date')
    AND t.created_at <= sqlc.arg('end_date')
    AND (sqlc.arg('categories')::INTEGER[] IS NULL OR t.subcategory_id = ANY(sqlc.arg('categories')::INTEGER[]))
GROUP BY cd.geo_location
ORDER BY count DESC;

-- name: GetHeatmapStats :one
SELECT
    COUNT(DISTINCT cd.geo_location) as total_locations,
    COUNT(*) as total_tickets,
    (COUNT(*)::FLOAT / NULLIF(COUNT(DISTINCT cd.geo_location), 0)::FLOAT) as avg_tickets_per_location
FROM complaint_details cd
INNER JOIN tickets t ON t.id = cd.ticket
WHERE 
    t.is_deleted = false 
    AND t.is_hidden = false;

-- name: GetTopProblemLocations :many
SELECT
    ST_Y(cd.geo_location::geometry)::DOUBLE PRECISION as lat,
    ST_X(cd.geo_location::geometry)::DOUBLE PRECISION as lng,
    COUNT(*) as ticket_count
FROM complaint_details cd
INNER JOIN tickets t ON t.id = cd.ticket
WHERE 
    t.is_deleted = false 
    AND t.is_hidden = false
GROUP BY cd.geo_location
ORDER BY ticket_count DESC
LIMIT 5;
