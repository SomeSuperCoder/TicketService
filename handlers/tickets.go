package handlers

import (
	"context"
	"fmt"

	"github.com/SomeSuperCoder/OnlineShop/internal/embeddings"
	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/google/uuid"
)

type TicketHandler struct {
	Repo *repository.Queries
}

// ==================== CREATE ====================

type PostTicketRequest struct {
	Body struct {
		Description   string  `json:"description"`
		SenderName    string  `json:"sender_name"`
		SenderPhone   string  `json:"sender_phone"`
		SenderEmail   string  `json:"sender_email"`
		Longitude     float64 `json:"longitude"`
		Latitude      float64 `json:"latitude"`
		SubcategoryID int32   `json:"subcategory_id"`
		DepartmentID  *int32  `json:"department_id,omitempty"`
	}
}

type PostTicketResponse struct {
	Body repository.Ticket
}

func (h *TicketHandler) Post(ctx context.Context, req *PostTicketRequest) (*PostTicketResponse, error) {
	resp := new(PostTicketResponse)

	// Generate embedding from description
	vector, err := embeddings.GetEmbedding(req.Body.Description)
	if err != nil {
		return nil, err
	}

	// Create GeoJSON or WKT point
	geoLocation := fmt.Sprintf("POINT(%f %f)", req.Body.Longitude, req.Body.Latitude)

	result, err := h.Repo.CreateTicketWithDefaults(ctx, repository.CreateTicketWithDefaultsParams{
		Description:   req.Body.Description,
		SenderName:    req.Body.SenderName,
		SenderPhone:   req.Body.SenderPhone,
		SenderEmail:   req.Body.SenderEmail,
		GeoLocation:   geoLocation,
		SubcategoryID: req.Body.SubcategoryID,
		DepartmentID:  req.Body.DepartmentID,
		Embedding:     vector,
	})
	if err != nil {
		return nil, err
	}

	resp.Body = result
	return resp, nil
}

// ==================== READ ====================

type GetTicketRequest struct {
	ID uuid.UUID `path:"id"`
}

type GetTicketResponse struct {
	Body repository.Ticket
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

type UpdateTicketRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Status        *repository.TicketStatus `json:"status,omitempty"`
		Complaints    []string                 `json:"complaints,omitempty"`
		Description   *string                  `json:"description,omitempty"`
		SubcategoryID *int32                   `json:"subcategory_id,omitempty"`
		DepartmentID  *int32                   `json:"department_id,omitempty"`
	}
}

type UpdateTicketResponse struct {
	Body repository.Ticket
}

func (h *TicketHandler) Update(ctx context.Context, req *UpdateTicketRequest) (*UpdateTicketResponse, error) {
	resp := new(UpdateTicketResponse)

	// First get the current ticket to check if we need to regenerate embedding
	currentTicket, err := h.Repo.GetTicket(ctx, repository.GetTicketParams{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}

	// Prepare update params
	params := repository.UpdateTicketParams{
		ID:            req.ID,
		Status:        currentTicket.Status, // Will be overridden by COALESCE
		Complaints:    currentTicket.Complaints,
		Description:   currentTicket.Description,
		SubcategoryID: currentTicket.SubcategoryID,
		DepartmentID:  currentTicket.DepartmentID,
		Embedding:     currentTicket.Embedding,
	}

	// Override with provided values
	if req.Body.Status != nil {
		params.Status = *req.Body.Status
	}
	if req.Body.Complaints != nil {
		params.Complaints = req.Body.Complaints
	}
	if req.Body.Description != nil {
		params.Description = *req.Body.Description
		// Regenerate embedding if description changed
		vector, err := embeddings.GetEmbedding(*req.Body.Description)
		if err != nil {
			return nil, err
		}
		params.Embedding = vector
	}
	if req.Body.SubcategoryID != nil {
		params.SubcategoryID = *req.Body.SubcategoryID
	}
	if req.Body.DepartmentID != nil {
		params.DepartmentID = req.Body.DepartmentID
	}

	result, err := h.Repo.UpdateTicket(ctx, params)
	if err != nil {
		return nil, err
	}

	resp.Body = result
	return resp, nil
}

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
