# Ticket History Feature

The history tracking system automatically records all changes made to tickets, including status changes, department assignments, tag modifications, comments, merges, and deletions.

## What's Tracked

The system tracks the following actions:
- `created` - When a ticket is first created
- `status_changed` - When ticket status changes (init → open → closed)
- `department_changed` - When a ticket is assigned to a different department
- `tags_added` - When tags are added to a ticket
- `tags_removed` - When tags are removed from a ticket
- `comment_added` - When a comment is added to a ticket
- `merged` - When duplicate tickets are merged
- `deleted` - When a ticket is soft-deleted

## Database Schema

The `ticket_history` table stores:
- `id` - Unique identifier for the history entry
- `ticket_id` - Reference to the ticket
- `action` - Type of change (enum)
- `old_value` - Previous state (JSONB)
- `new_value` - New state (JSONB)
- `user_id` - Optional UUID of the user who made the change
- `user_name` - Optional name of the user
- `user_email` - Optional email of the user
- `description` - Optional description of the change
- `created_at` - Timestamp of the change

## API Endpoints

### Get Ticket History
```
GET /api/v1/tickets/:id/history?limit=50&offset=0
```
Returns the complete history for a specific ticket, including all status changes and user actions.

Response:
```json
{
  "history": [
    {
      "id": "uuid",
      "ticket_id": "uuid",
      "action": "status_changed",
      "old_value": {"status": "init"},
      "new_value": {"status": "open"},
      "user_name": "John Doe",
      "user_email": "john@example.com",
      "created_at": "2024-01-01T12:00:00Z"
    }
  ],
  "total": 10
}
```

### Get Recent History (All Tickets)
```
GET /api/v1/history/recent?limit=50&offset=0
```
Returns recent history across all tickets.

### Add Comment
```
POST /api/v1/tickets/:id/comments
```
Adds a comment to a ticket and automatically creates a history entry.

Request:
```json
{
  "message": "Comment text",
  "user_name": "John Doe",
  "user_email": "john@example.com"
}
```

### Get Comments
```
GET /api/v1/tickets/:id/comments
```
Returns all comments for a ticket.

## Implementation

History entries are automatically created in the handlers:
- `TicketHandler.Post()` - Records ticket creation
- `TicketHandler.Update()` - Records status, department, and tag changes
- `TicketHandler.Delete()` - Records deletion
- `TicketHandler.Merge()` - Records merge operations
- `CommentsHandler.Post()` - Records comment additions

All history operations are wrapped in database transactions to ensure consistency.

## User Tracking

The system supports optional user tracking for all actions. When making changes, you can include:
- `user_id` - UUID of the authenticated user
- `user_name` - Display name
- `user_email` - Email address

This allows the frontend to display "who did what" in the history timeline.
