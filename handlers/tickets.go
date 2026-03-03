package handlers

import (
	"context"
	"fmt"

	"github.com/SomeSuperCoder/OnlineShop/internal/embeddings"
	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TicketHandler struct {
	Repo *repository.Queries
	Pool *pgxpool.Pool
}

// ==================== CREATE ====================

type PostTicketRequest struct {
	Body struct {
		Description   string  `json:"description"`
		SenderName    string  `json:"sender_name"`
		SenderPhone   string  `json:"sender_phone" default:"+1 500 555 0006"`
		SenderEmail   *string `json:"sender_email" default:"test@example.com"`
		Longitude     float64 `json:"longitude"`
		Latitude      float64 `json:"latitude"`
		SubcategoryID int32   `json:"subcategory_id"`
		DepartmentID  *int32  `json:"department_id,omitempty"`
	}
}

type PostTicketResponse struct {
	Body struct {
		Ticket           repository.Ticket          `json:"ticket"`
		ComplaintDetails repository.ComplaintDetail `json:"complaint_details"`
	}
}

func (h *TicketHandler) Post(ctx context.Context, req *PostTicketRequest) (*PostTicketResponse, error) {
	resp := new(PostTicketResponse)

	tx, err := h.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := h.Repo.WithTx(tx)

	// Generate embedding from description
	vector, err := embeddings.GetEmbedding(req.Body.Description)
	if err != nil {
		return nil, err
	}

	// Create GeoJSON or WKT point
	geoLocation := fmt.Sprintf("POINT(%f %f)", req.Body.Longitude, req.Body.Latitude)

	result, err := qtx.CreateTicketWithDefaults(ctx, repository.CreateTicketWithDefaultsParams{
		Description:   req.Body.Description,
		SubcategoryID: req.Body.SubcategoryID,
		DepartmentID:  req.Body.DepartmentID,
		Embedding:     vector,
	})
	if err != nil {
		return nil, err
	}
	details, err := qtx.CreateComplaint(ctx, repository.CreateComplaintParams{
		Ticket:      result.ID,
		Description: req.Body.Description,
		SenderName:  req.Body.SenderName,
		SenderPhone: &req.Body.SenderPhone,
		SenderEmail: req.Body.SenderEmail,
		GeoLocation: geoLocation,
	})
	if err != nil {
		return nil, err
	}

	resp.Body.Ticket = result
	resp.Body.ComplaintDetails = details

	return resp, nil
}

// ==================== READ ====================

type GetTicketRequest struct {
	ID uuid.UUID `path:"id"`
}

type GetTicketResponse struct {
	Body repository.GetTicketRow
}

func (h *TicketHandler) Get(ctx context.Context, req *GetTicketRequest) (*GetTicketResponse, error) {
	resp := new(GetTicketResponse)

	result, err := h.Repo.GetTicket(ctx, repository.GetTicketParams{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}

	resp.Body = result
	return resp, nil
}

type ListTicketsRequest struct {
	Limit  int32 `query:"limit" default:"10" maximum:"100"`
	Offset int32 `query:"offset" default:"0"`
}

type ListTicketsResponse struct {
	Body struct {
		Tickets []repository.Ticket `json:"tickets"`
		Total   int64               `json:"total"`
	}
}

func (h *TicketHandler) List(ctx context.Context, req *ListTicketsRequest) (*ListTicketsResponse, error) {
	resp := new(ListTicketsResponse)

	tickets, err := h.Repo.ListTickets(ctx, repository.ListTicketsParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return nil, err
	}

	total, err := h.Repo.CountTickets(ctx)
	if err != nil {
		return nil, err
	}

	resp.Body.Tickets = tickets
	resp.Body.Total = total
	return resp, nil
}

// ==================== SEARCH ====================

type SearchByMeaningRequest struct {
	Query string `query:"query"`
	Limit int32  `query:"limit" default:"10" maximum:"50"`
}

type SearchByMeaningResponse struct {
	Body struct {
		Tickets []repository.Ticket `json:"tickets"`
	}
}

func (h *TicketHandler) SearchByMeaning(ctx context.Context, req *SearchByMeaningRequest) (*SearchByMeaningResponse, error) {
	resp := new(SearchByMeaningResponse)

	vector, err := embeddings.GetEmbedding(req.Query)
	if err != nil {
		return nil, err
	}

	result, err := h.Repo.SearchTicketsByEmbedding(ctx, repository.SearchTicketsByEmbeddingParams{
		Embedding: vector,
		Limit:     req.Limit,
	})
	if err != nil {
		return nil, err
	}

	resp.Body.Tickets = result
	return resp, nil
}

// ==================== UPDATE ====================
// ==================== DELETE / HIDE ====================
type DeleteTicketRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeleteTicketResponse struct {
	Body struct {
		Message string `json:"message"`
	}
}

func (h *TicketHandler) Delete(ctx context.Context, req *DeleteTicketRequest) (*DeleteTicketResponse, error) {
	resp := new(DeleteTicketResponse)

	err := h.Repo.DeleteTicket(ctx, repository.DeleteTicketParams{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}

	resp.Body.Message = "Ticket deleted successfully"
	return resp, nil
}

// ==================== COUNT ====================

type CountTicketsRequest struct{}

type CountTicketsResponse struct {
	Body struct {
		Count int64 `json:"count"`
	}
}

func (h *TicketHandler) Count(ctx context.Context, req *CountTicketsRequest) (*CountTicketsResponse, error) {
	resp := new(CountTicketsResponse)

	count, err := h.Repo.CountTickets(ctx)
	if err != nil {
		return nil, err
	}

	resp.Body.Count = count
	return resp, nil
}
