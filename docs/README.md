# API Documentation

Complete API reference for frontend integration.

**Base URL**: `http://localhost:8888/api/v1`

**Swagger UI**: http://localhost:8888/api/v1/docs

## Documentation Files

- [API_REFERENCE.md](API_REFERENCE.md) - Complete API endpoint reference
- [ADMIN_USERS_API.md](ADMIN_USERS_API.md) - User management endpoints
- [KPI_ENDPOINT.md](KPI_ENDPOINT.md) - KPI monitoring endpoint

## Quick Reference

### Public Endpoints
- `GET /public/categories` - Get category tree
- `POST /public/tickets` - Submit new ticket

### Tickets
- `GET /tickets` - List tickets (with filters and semantic search)
- `GET /tickets/{id}` - Get ticket details
- `PATCH /tickets/{id}` - Update ticket
- `DELETE /tickets/{id}` - Delete ticket
- `POST /tickets/merge` - Merge duplicates

### Comments
- `POST /tickets/{id}/comments` - Add comment
- `GET /tickets/{id}/comments` - Get comments

### Statistics
- `GET /statistics/summary` - Overall stats
- `GET /statistics/categories` - Category breakdown
- `GET /statistics/dynamics` - Time series data

### Monitoring
- `GET /monitoring/kpi` - Key performance indicators
- `GET /monitoring/departments` - Department efficiency
- `GET /monitoring/overdue` - Overdue tickets

### Heatmap
- `GET /heatmap/points` - Geographic points
- `GET /heatmap/stats` - Heatmap statistics

### Admin
- `GET /admin/users` - List users
- `POST /admin/users` - Create user
- `PATCH /admin/users/{id}` - Update user
- `DELETE /admin/users/{id}` - Delete user

### History
- `GET /tickets/{id}/history` - Ticket history
- `GET /history/recent` - Recent activity

## Data Models

### Ticket Status
- `init` - Initial state
- `open` - Being processed
- `closed` - Resolved

### User Roles
- `admin` - Full access
- `org` - Organization/department level
- `executor` - Worker level

### User Status
- `active` - Can access system
- `blocked` - Access denied

## Common Patterns

### Pagination
```
?limit=20&offset=0
```

### Filtering
```
?status_id=open&subcategory_id=1
```

### Semantic Search
```
?query=проблема+с+отоплением
```

### Date Ranges
```
?start_date=2026-01-01&end_date=2026-03-05
```

## Response Format

All responses are JSON. Errors include:
```json
{
  "error": "Error message",
  "details": "Additional details"
}
```

## Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `404` - Not Found
- `500` - Internal Server Error
