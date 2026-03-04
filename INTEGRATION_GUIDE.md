# History & Comments Integration Guide

## Overview

The history tracking system is now fully implemented according to the llms.txt specification. It tracks all ticket changes and user actions.

## New Handlers

### 1. HistoryHandler (`handlers/history.go`)
Provides endpoints to retrieve ticket history.

### 2. CommentsHandler (`handlers/comments.go`)
Handles comment creation and retrieval, automatically creating history entries.

## API Routes to Register

Consistent with your project's "tickets" terminology:

```go
// Comments
POST   /api/v1/tickets/:id/comments     -> CommentsHandler.Post
GET    /api/v1/tickets/:id/comments     -> CommentsHandler.Get

// History
GET    /api/v1/tickets/:id/history      -> HistoryHandler.GetTicketHistory

// Optional: Recent history across all tickets
GET    /api/v1/history/recent           -> HistoryHandler.GetRecentHistory
```

## Example Integration with Huma

```go
package main

import (
    "github.com/SomeSuperCoder/OnlineShop/handlers"
    "github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API, ticketHandler *handlers.TicketHandler, 
                    historyHandler *handlers.HistoryHandler,
                    commentsHandler *handlers.CommentsHandler) {
    
    // Comments endpoints
    huma.Register(api, huma.Operation{
        OperationID: "post-comment",
        Method:      "POST",
        Path:        "/tickets/{id}/comments",
        Summary:     "Add a comment to a ticket",
    }, commentsHandler.Post)
    
    huma.Register(api, huma.Operation{
        OperationID: "get-comments",
        Method:      "GET",
        Path:        "/tickets/{id}/comments",
        Summary:     "Get all comments for a ticket",
    }, commentsHandler.Get)
    
    // History endpoint
    huma.Register(api, huma.Operation{
        OperationID: "get-ticket-history",
        Method:      "GET",
        Path:        "/tickets/{id}/history",
        Summary:     "Get history of actions for a ticket",
    }, historyHandler.GetTicketHistory)
    
    // Optional: Recent history
    huma.Register(api, huma.Operation{
        OperationID: "get-recent-history",
        Method:      "GET",
        Path:        "/history/recent",
        Summary:     "Get recent history across all tickets",
    }, historyHandler.GetRecentHistory)
}
```

## Database Migration

Run the migration to create the history table:

```bash
goose -dir db/migrations postgres "your-connection-string" up
```

## Automatic History Tracking

History is automatically recorded for:
- ✅ Ticket creation
- ✅ Status changes
- ✅ Department changes
- ✅ Tag additions/removals
- ✅ Comments
- ✅ Merges
- ✅ Deletions

No additional code needed - it's built into the existing handlers!

## User Context (Optional)

To track which user made changes, update the handlers to accept user context:

```go
// Example: Add user info to comment
req := &PostCommentRequest{
    TicketID: ticketID,
    Body: struct {
        Message   string  `json:"message"`
        UserName  *string `json:"user_name,omitempty"`
        UserEmail *string `json:"user_email,omitempty"`
    }{
        Message:   "This is a comment",
        UserName:  &userName,
        UserEmail: &userEmail,
    },
}
```

You can extract user info from JWT tokens and pass it to the handlers.
