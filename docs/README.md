# API Documentation

Complete API reference for frontend integration.

**Base URL**: `http://localhost:8888/api/v1`

**Swagger UI**: http://localhost:8888/api/v1/docs

## Getting Started

```bash
# Start databases
make databases

# Run migrations
make migrate

# Load mock data (recommended for testing)
make mock

# Start API server
make serve
```

The mock data script (`mock/main.go`) populates the database with realistic, diverse test data including:

**Dictionaries:**
- 8 categories (ЖКХ, Дороги, Благоустройство, Безопасность, Связь, Здравоохранение, Образование, Экология)
- 30+ subcategories
- 8 departments
- 10 users (admin, orgs, executors with different statuses)
- 7 tags
- 6 sources

**Tickets (25 diverse scenarios):**
- 6 **closed** tickets (resolved, 15-45 days old)
- 6 **open** tickets (in progress, 1-6 days old)
- 6 **overdue** tickets (open for 8-20 days, ensuring overdue_count > 0)
- 5 **init** tickets (new, unassigned, recent)

**Features:**
- ✅ **Real AI embeddings** - Generated using Ollama (nomic-embed-text-v2-moe model)
- ✅ **Semantic similarity** - Tickets can be found using natural language search
- ✅ **Geographic diversity** - Different locations across Cheboksary for heatmap
- ✅ **Time series data** - Tickets spread across 45 days for dynamics charts
- ✅ **Rich metadata** - Comments, history, tags, department assignments
- ✅ **Realistic scenarios** - Heating problems, road damage, stray dogs, internet issues, etc.

All API endpoints will return meaningful, non-default values with this data loaded.

**Requirements:**
- Ollama running locally: `ollama run nomic-embed-text-v2-moe`
- If Ollama is not available, fallback embeddings will be used (semantic search won't work as well)

## Documentation Files

- [API_REFERENCE.md](API_REFERENCE.md) - Complete API endpoint reference with examples

For detailed endpoint documentation, see API_REFERENCE.md or use the interactive Swagger UI.

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
