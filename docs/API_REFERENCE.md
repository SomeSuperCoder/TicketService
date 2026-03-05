# API Reference

Complete API reference for the Intelligent Citizen Appeals Processing System.

Base URL: `http://localhost:8888/api/v1`

## Table of Contents

- [Public Endpoints](#public-endpoints)
- [Ticket Management](#ticket-management)
- [Comments](#comments)
- [Statistics](#statistics)
- [Monitoring](#monitoring)
- [Heatmap](#heatmap)
- [Admin - User Management](#admin---user-management)
- [History](#history)

---

## Public Endpoints

### Get Categories

Get a tree of categories and subcategories for the appeal form.

**Endpoint**: `GET /public/categories`

**Response**:
```json
{
  "categories": [
    {
      "category_info": {
        "id": 1,
        "name": "ЖКХ"
      },
      "subcategories": [
        {
          "id": 1,
          "category_id": 1,
          "name": "Плата за услуги ЖКУ"
        },
        {
          "id": 2,
          "category_id": 1,
          "name": "Качество услуг"
        }
      ]
    }
  ]
}
```

### Create Ticket (Public)

Submit a new citizen appeal.

**Endpoint**: `POST /public/tickets`

**Request Body**:
```json
{
  "description": "Проблема с отоплением в доме",
  "sender_name": "Иванов Иван Иванович",
  "sender_phone": "+7 900 123 45 67",
  "sender_email": "ivanov@example.com",
  "longitude": 47.2501,
  "latitude": 56.1324,
  "subcategory_id": 1,
  "department_id": 2
}
```

**Required Fields**:
- `description` - Appeal text
- `sender_name` - Full name or nickname
- `longitude`, `latitude` - Geographic coordinates
- `subcategory_id` - Subcategory ID from categories endpoint

**Optional Fields**:
- `sender_phone` - Phone number
- `sender_email` - Email address
- `department_id` - Department to assign to

**Note**: At least one contact method (phone or email) must be provided unless anonymous.

**Response**:
```json
{
  "ticket": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "init",
    "description": "Проблема с отоплением в доме",
    "subcategory_id": 1,
    "department_id": 2,
    "is_hidden": false,
    "is_deleted": false,
    "created_at": "2026-03-05T10:30:00Z"
  },
  "complaint_details": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "ticket": "550e8400-e29b-41d4-a716-446655440000",
    "description": "Проблема с отоплением в доме",
    "sender_name": "Иванов Иван Иванович",
    "sender_phone": "+7 900 123 45 67",
    "sender_email": "ivanov@example.com",
    "geo_location": "POINT(47.2501 56.1324)"
  }
}
```

---

## Ticket Management

### List Tickets

Get a paginated list of tickets with filtering.

**Endpoint**: `GET /tickets`

**Query Parameters**:
- `query` (optional) - Semantic search query (uses AI embeddings)
- `limit` (optional) - Results per page (default: 10, max: 100)
- `offset` (optional) - Pagination offset (default: 0)
- `status_id` (optional) - Filter by status: `init`, `open`, `closed`
- `subcategory_id` (optional) - Filter by subcategory ID

**Example**:
```
GET /tickets?limit=20&offset=0&status_id=open&subcategory_id=1
```

**Response**:
```json
{
  "tickets": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "status": "open",
      "description": "Проблема с отоплением",
      "subcategory_id": 1,
      "subcategory_name": "Плата за услуги ЖКУ",
      "category_name": "ЖКХ",
      "department_id": 2,
      "department_name": "Отдел ЖКХ",
      "created_at": "2026-03-05T10:30:00Z",
      "similarity": 0.95
    }
  ],
  "total": 150
}
```

### Get Ticket

Get detailed information about a specific ticket.

**Endpoint**: `GET /tickets/{id}`

**Path Parameters**:
- `id` - Ticket UUID

**Response**:
```json
{
  "ticket": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "open",
    "description": "Проблема с отоплением в доме",
    "subcategory_id": 1,
    "subcategory_name": "Плата за услуги ЖКУ",
    "category_name": "ЖКХ",
    "department_id": 2,
    "department_name": "Отдел ЖКХ",
    "created_at": "2026-03-05T10:30:00Z",
    "tags": [1, 3]
  },
  "details": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "ticket": "550e8400-e29b-41d4-a716-446655440000",
      "description": "Проблема с отоплением в доме",
      "sender_name": "Иванов Иван Иванович",
      "sender_phone": "+7 900 123 45 67",
      "sender_email": "ivanov@example.com",
      "latitude": 56.1324,
      "longitude": 47.2501
    }
  ]
}
```

### Update Ticket

Update ticket status, department, or tags.

**Endpoint**: `PATCH /tickets/{id}`

**Path Parameters**:
- `id` - Ticket UUID

**Request Body** (all fields optional):
```json
{
  "status": "open",
  "department_id": 3,
  "add_tags": [1, 2],
  "remove_tags": [3]
}
```

**Fields**:
- `status` - New status: `init`, `open`, `closed`
- `department_id` - Assign to different department
- `add_tags` - Array of tag IDs to add
- `remove_tags` - Array of tag IDs to remove

**Response**:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "open",
  "description": "Проблема с отоплением в доме",
  "subcategory_id": 1,
  "department_id": 3,
  "is_hidden": false,
  "is_deleted": false,
  "created_at": "2026-03-05T10:30:00Z"
}
```

### Delete Ticket

Permanently delete a ticket (soft delete).

**Endpoint**: `DELETE /tickets/{id}`

**Path Parameters**:
- `id` - Ticket UUID

**Response**:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "open",
  "description": "Проблема с отоплением в доме",
  "subcategory_id": 1,
  "department_id": 2,
  "is_hidden": false,
  "is_deleted": true,
  "created_at": "2026-03-05T10:30:00Z"
}
```

### Merge Duplicate Tickets

Merge duplicate tickets into one original ticket.

**Endpoint**: `POST /tickets/merge`

**Request Body**:
```json
{
  "original": "550e8400-e29b-41d4-a716-446655440000",
  "duplicates": [
    "660e8400-e29b-41d4-a716-446655440001",
    "770e8400-e29b-41d4-a716-446655440002"
  ]
}
```

**Process**:
1. All complaint details from duplicates are transferred to the original ticket
2. Embeddings are averaged
3. Duplicate tickets are deleted
4. History entry is created

**Response**: `{}`

---

## Comments

### Add Comment

Add a comment to a ticket.

**Endpoint**: `POST /tickets/{id}/comments`

**Path Parameters**:
- `id` - Ticket UUID

**Request Body**:
```json
{
  "message": "Проблема решена, отопление восстановлено",
  "user_name": "Петров П.П.",
  "user_email": "petrov@admin.local"
}
```

**Response**:
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440003",
  "ticket": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Проблема решена, отопление восстановлено"
}
```

### Get Comments

Get all comments for a ticket.

**Endpoint**: `GET /tickets/{id}/comments`

**Path Parameters**:
- `id` - Ticket UUID

**Response**:
```json
{
  "comments": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440003",
      "ticket": "550e8400-e29b-41d4-a716-446655440000",
      "message": "Проблема решена, отопление восстановлено"
    }
  ]
}
```

---

## Statistics

### Get Summary Statistics

Get overall system statistics.

**Endpoint**: `GET /statistics/summary`

**Response**:
```json
{
  "total_tickets": 1250,
  "open_tickets": 320,
  "closed_tickets": 850,
  "init_tickets": 80,
  "active_executors": 45
}
```

### Get Category Statistics

Get ticket counts by category.

**Endpoint**: `GET /statistics/categories`

**Response**:
```json
[
  {
    "category_id": 1,
    "category_name": "ЖКХ",
    "ticket_count": 450
  },
  {
    "category_id": 2,
    "category_name": "Дороги",
    "ticket_count": 320
  }
]
```

### Get Ticket Dynamics

Get ticket dynamics over time (received and resolved).

**Endpoint**: `GET /statistics/dynamics`

**Query Parameters**:
- `period` (optional) - Grouping period: `day` (default), `week`
- `start_date` (optional) - Start date (YYYY-MM-DD), default: 30 days ago
- `end_date` (optional) - End date (YYYY-MM-DD), default: today

**Example**:
```
GET /statistics/dynamics?period=week&start_date=2026-01-01&end_date=2026-03-05
```

**Response**:
```json
{
  "period": "week",
  "data": [
    {
      "date": "2026-02-03",
      "received": 45,
      "resolved": 38
    },
    {
      "date": "2026-02-10",
      "received": 52,
      "resolved": 41
    }
  ]
}
```

---

## Monitoring

### Get KPI Metrics

Get key performance indicators.

**Endpoint**: `GET /monitoring/kpi`

**Query Parameters**:
- `period` (optional) - Time period: `week`, `month` (default), `year`

**Response**:
```json
{
  "avg_response_days": 2.5,
  "overdue_count": 15,
  "satisfaction_index": 75.5
}
```

See [KPI_ENDPOINT.md](KPI_ENDPOINT.md) for detailed metric explanations.

### Get Department Efficiency

Get efficiency metrics for all departments.

**Endpoint**: `GET /monitoring/departments`

**Response**:
```json
[
  {
    "department_id": 1,
    "department_name": "Отдел ЖКХ",
    "total_in_progress": 45,
    "overdue_count": 5,
    "avg_resolution_days": 3.2,
    "trend": "improving"
  }
]
```

### Get Overdue Tickets

Get list of overdue tickets with lost days metric.

**Endpoint**: `GET /monitoring/overdue`

**Response**:
```json
[
  {
    "ticket_id": "550e8400-e29b-41d4-a716-446655440000",
    "description": "Проблема с отоплением",
    "status": "open",
    "department_name": "Отдел ЖКХ",
    "lost_days": 3,
    "created_at": "2026-02-20T10:30:00Z"
  }
]
```

---

## Heatmap

### Get Heatmap Points

Get geographic points for heatmap visualization.

**Endpoint**: `GET /heatmap/points`

**Query Parameters**:
- `period` (optional) - Time period: `week`, `month` (default), `year`
- `categories` (optional) - Array of category IDs to filter

**Example**:
```
GET /heatmap/points?period=month&categories=1&categories=2
```

**Response**:
```json
{
  "points": [
    {
      "lat": 56.1324,
      "lng": 47.2501,
      "intensity": 15,
      "count": 15
    }
  ]
}
```

### Get Heatmap Statistics

Get statistics for the heatmap.

**Endpoint**: `GET /heatmap/stats`

**Response**:
```json
{
  "top_locations": [
    {
      "lat": 56.1324,
      "lng": 47.2501,
      "ticket_count": 25
    }
  ],
  "total_locations": 150,
  "avg_tickets_per_location": 3.5
}
```

---

## Admin - User Management

See [ADMIN_USERS_API.md](ADMIN_USERS_API.md) for complete documentation.

### List Users

**Endpoint**: `GET /admin/users`

### Create User

**Endpoint**: `POST /admin/users`

### Update User

**Endpoint**: `PATCH /admin/users/{id}`

### Delete User

**Endpoint**: `DELETE /admin/users/{id}`

---

## History

### Get Ticket History

Get complete history of actions for a ticket.

**Endpoint**: `GET /tickets/{id}/history`

**Path Parameters**:
- `id` - Ticket UUID

**Response**:
```json
[
  {
    "id": "990e8400-e29b-41d4-a716-446655440004",
    "ticket_id": "550e8400-e29b-41d4-a716-446655440000",
    "action": "created",
    "new_value": {
      "status": "init",
      "subcategory_id": 1
    },
    "created_at": "2026-03-05T10:30:00Z"
  },
  {
    "id": "aa0e8400-e29b-41d4-a716-446655440005",
    "ticket_id": "550e8400-e29b-41d4-a716-446655440000",
    "action": "status_changed",
    "old_value": {
      "status": "init"
    },
    "new_value": {
      "status": "open"
    },
    "user_name": "Петров П.П.",
    "created_at": "2026-03-05T11:00:00Z"
  }
]
```

**History Actions**:
- `created` - Ticket created
- `status_changed` - Status updated
- `department_changed` - Department reassigned
- `tags_added` - Tags added
- `tags_removed` - Tags removed
- `comment_added` - Comment added
- `merged` - Tickets merged
- `deleted` - Ticket deleted

### Get Recent History

Get recent history across all tickets.

**Endpoint**: `GET /history/recent`

**Query Parameters**:
- `limit` (optional) - Number of entries (default: 50)

**Response**: Same format as ticket history.

---

## Error Responses

All endpoints return standard error responses:

```json
{
  "error": "Error message",
  "details": "Additional error details"
}
```

**HTTP Status Codes**:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `404` - Not Found
- `500` - Internal Server Error

---

## Authentication

Currently, the API does not require authentication. JWT authentication is configured but not enforced. Future versions will require JWT tokens for CRM endpoints.

**Planned Authentication Header**:
```
Authorization: Bearer <jwt_token>
```
