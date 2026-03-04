package handlers

import (
	"context"
	"encoding/json"

	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentsHandler struct {
	Repo *repository.Queries
	Pool *pgxpool.Pool
}

// ==================== CREATE COMMENT ====================

type PostCommentRequest struct {
	TicketID uuid.UUID `path:"id"`
	Body     struct {
		Message   string  `json:"message"`
		UserName  *string `json:"user_name,omitempty"`
		UserEmail *string `json:"user_email,omitempty"`
	}
}

type PostCommentResponse struct {
	Body repository.Comment `json:"comment"`
}

func (h *CommentsHandler) Post(ctx context.Context, req *PostCommentRequest) (*PostCommentResponse, error) {
	resp := new(PostCommentResponse)

	tx, err := h.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := h.Repo.WithTx(tx)

	// Create the comment
	comment, err := qtx.CreateComment(ctx, repository.CreateCommentParams{
		Ticket:  req.TicketID,
		Message: req.Body.Message,
	})
	if err != nil {
		return nil, err
	}

	// Record history entry for the comment
	newValue, _ := json.Marshal(map[string]interface{}{
		"comment_id": comment.ID,
		"message":    comment.Message,
	})
	_, err = qtx.CreateHistoryEntry(ctx, repository.CreateHistoryEntryParams{
		TicketID:  req.TicketID,
		Action:    repository.HistoryActionCommentAdded,
		NewValue:  newValue,
		UserName:  req.Body.UserName,
		UserEmail: req.Body.UserEmail,
	})
	if err != nil {
		return nil, err
	}

	resp.Body = comment

	tx.Commit(ctx)

	return resp, nil
}

// ==================== GET COMMENTS ====================

type GetCommentsRequest struct {
	TicketID uuid.UUID `path:"id"`
}

type GetCommentsResponse struct {
	Body struct {
		Comments []repository.Comment `json:"comments"`
	}
}

func (h *CommentsHandler) Get(ctx context.Context, req *GetCommentsRequest) (*GetCommentsResponse, error) {
	resp := new(GetCommentsResponse)

	comments, err := h.Repo.GetCommentsForTicket(ctx, repository.GetCommentsForTicketParams{
		Ticket: req.TicketID,
	})
	if err != nil {
		return nil, err
	}

	resp.Body.Comments = comments

	return resp, nil
}
