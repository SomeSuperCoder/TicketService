package handlers

import (
	"context"

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
		Description   string   `json:"description"`
		Complaints    []string `json:"complaints"`
		SubcategoryID int32    `json:"subcategory_id"`
		DepartmentID  *int32   `json:"department_id,omitempty"`
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

	result, err := h.Repo.CreateTicketWithDefaults(ctx, repository.CreateTicketWithDefaultsParams{
		Description:   req.Body.Description,
		Complaints:    req.Body.Complaints,
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
	ID uuid.UUID `param:"id"`
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

type GetTicketsByStatusRequest struct {
	Status repository.TicketStatus `query:"status" enum:"init,open,closed"`
	Limit  int32                   `query:"limit" default:"10" maximum:"100"`
	Offset int32                   `query:"offset" default:"0"`
}

type GetTicketsByStatusResponse struct {
	Body struct {
		Tickets []repository.Ticket `json:"tickets"`
		Total   int64               `json:"total"`
	}
}

func (h *TicketHandler) GetByStatus(ctx context.Context, req *GetTicketsByStatusRequest) (*GetTicketsByStatusResponse, error) {
	resp := new(GetTicketsByStatusResponse)

	tickets, err := h.Repo.GetTicketsByStatus(ctx, repository.GetTicketsByStatusParams{
		Status: req.Status,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return nil, err
	}

	total, err := h.Repo.CountTicketsByStatus(ctx, repository.CountTicketsByStatusParams{
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	resp.Body.Tickets = tickets
	resp.Body.Total = total
	return resp, nil
}

type GetTicketsBySubcategoryRequest struct {
	SubcategoryID int32 `param:"subcategoryId"`
	Limit         int32 `query:"limit" default:"10" maximum:"100"`
	Offset        int32 `query:"offset" default:"0"`
}

type GetTicketsBySubcategoryResponse struct {
	Body struct {
		Tickets []repository.Ticket `json:"tickets"`
		Total   int64               `json:"total"`
	}
}

func (h *TicketHandler) GetBySubcategory(ctx context.Context, req *GetTicketsBySubcategoryRequest) (*GetTicketsBySubcategoryResponse, error) {
	resp := new(GetTicketsBySubcategoryResponse)

	tickets, err := h.Repo.GetTicketsBySubcategory(ctx, repository.GetTicketsBySubcategoryParams{
		SubcategoryID: req.SubcategoryID,
		Limit:         req.Limit,
		Offset:        req.Offset,
	})
	if err != nil {
		return nil, err
	}

	resp.Body.Tickets = tickets
	// Note: You might want to add a CountTicketsBySubcategory query if needed
	return resp, nil
}

type GetTicketsByDepartmentRequest struct {
	DepartmentID int32 `param:"departmentId"`
	Limit        int32 `query:"limit" default:"10" maximum:"100"`
	Offset       int32 `query:"offset" default:"0"`
}

type GetTicketsByDepartmentResponse struct {
	Body struct {
		Tickets []repository.Ticket `json:"tickets"`
	}
}

func (h *TicketHandler) GetByDepartment(ctx context.Context, req *GetTicketsByDepartmentRequest) (*GetTicketsByDepartmentResponse, error) {
	resp := new(GetTicketsByDepartmentResponse)

	tickets, err := h.Repo.GetTicketsByDepartment(ctx, repository.GetTicketsByDepartmentParams{
		DepartmentID: &req.DepartmentID,
		Limit:        req.Limit,
		Offset:       req.Offset,
	})
	if err != nil {
		return nil, err
	}

	resp.Body.Tickets = tickets
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
	ID   uuid.UUID `param:"id"`
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

type UpdateTicketStatusRequest struct {
	ID   uuid.UUID `param:"id"`
	Body struct {
		Status repository.TicketStatus `json:"status" enum:"init,open,closed"`
	}
}

type UpdateTicketStatusResponse struct {
	Body repository.Ticket
}

func (h *TicketHandler) UpdateStatus(ctx context.Context, req *UpdateTicketStatusRequest) (*UpdateTicketStatusResponse, error) {
	resp := new(UpdateTicketStatusResponse)

	result, err := h.Repo.UpdateTicketStatus(ctx, repository.UpdateTicketStatusParams{
		ID:     req.ID,
		Status: req.Body.Status,
	})
	if err != nil {
		return nil, err
	}

	resp.Body = result
	return resp, nil
}

// ==================== DELETE / HIDE ====================

type HideTicketRequest struct {
	ID uuid.UUID `param:"id"`
}

type HideTicketResponse struct {
	Body struct {
		Message string `json:"message"`
	}
}

func (h *TicketHandler) Hide(ctx context.Context, req *HideTicketRequest) (*HideTicketResponse, error) {
	resp := new(HideTicketResponse)

	err := h.Repo.HideTicket(ctx, repository.HideTicketParams{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}

	resp.Body.Message = "Ticket hidden successfully"
	return resp, nil
}

type DeleteTicketRequest struct {
	ID uuid.UUID `param:"id"`
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
