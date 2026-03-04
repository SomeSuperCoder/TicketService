package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HistoryHandler struct {
	Repo *repository.Queries
	Pool *pgxpool.Pool
}

// ==================== GET TICKET HISTORY ====================

type GetTicketHistoryRequest struct {
	TicketID uuid.UUID `path:"id"`
	Limit    int32     `query:"limit" default:"50" maximum:"200"`
	Offset   int32     `query:"offset" default:"0"`
}

type GetTicketHistoryResponse struct {
	Body struct {
		History []repository.TicketHistory `json:"history"`
		Total   int64                      `json:"total"`
	}
}

func (h *HistoryHandler) GetTicketHistory(ctx context.Context, req *GetTicketHistoryRequest) (*GetTicketHistoryResponse, error) {
	resp := new(GetTicketHistoryResponse)

	history, err := h.Repo.GetTicketHistory(ctx, repository.GetTicketHistoryParams{
		TicketID: req.TicketID,
		Limit:    req.Limit,
		Offset:   req.Offset,
	})
	if err != nil {
		return nil, err
	}

	total, err := h.Repo.GetTicketHistoryCount(ctx, repository.GetTicketHistoryCountParams{
		TicketID: req.TicketID,
	})
	if err != nil {
		return nil, err
	}

	resp.Body.History = history
	resp.Body.Total = total

	return resp, nil
}

// ==================== GET RECENT HISTORY ====================

type GetRecentHistoryRequest struct {
	Limit  int32 `query:"limit" default:"50" maximum:"200"`
	Offset int32 `query:"offset" default:"0"`
}

type GetRecentHistoryResponse struct {
	Body struct {
		History []repository.GetRecentHistoryRow `json:"history"`
	}
}

func (h *HistoryHandler) GetRecentHistory(ctx context.Context, req *GetRecentHistoryRequest) (*GetRecentHistoryResponse, error) {
	resp := new(GetRecentHistoryResponse)

	history, err := h.Repo.GetRecentHistory(ctx, repository.GetRecentHistoryParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return nil, err
	}

	resp.Body.History = history

	return resp, nil
}
