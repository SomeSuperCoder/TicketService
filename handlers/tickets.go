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

	tx.Commit(ctx)

	return resp, nil
}

// ==================== READ ====================

type GetTicketRequest struct {
	ID uuid.UUID `path:"id"`
}

type GetTicketResponse struct {
	Body struct {
		Ticket  repository.GetTicketRow             `json:"ticket"`
		Details []repository.GetDetailsForTicketRow `json:"details"`
	}
}

func (h *TicketHandler) Get(ctx context.Context, req *GetTicketRequest) (*GetTicketResponse, error) {
	resp := new(GetTicketResponse)

	tx, err := h.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := h.Repo.WithTx(tx)

	result, err := qtx.GetTicket(ctx, repository.GetTicketParams{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}

	details, err := qtx.GetDetailsForTicket(ctx, repository.GetDetailsForTicketParams{
		Ticket: req.ID,
	})
	if err != nil {
		return nil, err
	}

	resp.Body.Ticket = result
	resp.Body.Details = details

	tx.Commit(ctx)

	return resp, nil
}

// Left: `tags`
type ListTicketsRequest struct {
	Query         string                  `query:"query"`
	Limit         int32                   `query:"limit" default:"10" maximum:"100"`
	Offset        int32                   `query:"offset" default:"0"`
	Status        repository.TicketStatus `query:"status_id"`
	SubcategoryID int32                   `query:"subcategory_id"`
}

type ListTicketsResponse struct {
	Body struct {
		Tickets []repository.ListTicketsRow `json:"tickets"`
		Total   int64                       `json:"total"`
	}
}

func (h *TicketHandler) List(ctx context.Context, req *ListTicketsRequest) (*ListTicketsResponse, error) {
	resp := new(ListTicketsResponse)

	vector, err := embeddings.GetEmbedding(req.Query)
	if err != nil {
		return nil, err
	}

	// Defaults
	var statusValid bool = false
	if req.Status != repository.TicketStatusNone {
		statusValid = true
	}
	var subcategoryID *int32
	if req.SubcategoryID != 0 {
		subcategoryID = &req.SubcategoryID
	}

	tickets, err := h.Repo.ListTickets(ctx, repository.ListTicketsParams{
		Limit:  req.Limit,
		Offset: req.Offset,
		Status: repository.NullTicketStatus{
			TicketStatus: req.Status,
			Valid:        statusValid,
		},
		Subcategory: subcategoryID,
		Embedding:   vector,
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

// ==================== UPDATE ====================
type UpdateTicketRequest struct {
	TicketID uuid.UUID `path:"id"`
	Body     struct {
		Status       *repository.TicketStatus `json:"status,omitempty"`
		DepartmentID *int32                   `json:"department_id,omitempty"`
		AddTags      []int32                  `json:"add_tags,omitempty"`
		RemoveTags   []int32                  `json:"remove_tags,omitempty"`
	}
}
type UpdateTicketResponse struct {
	Body repository.UpdateTicketSimpleRow
}

func (h *TicketHandler) Update(ctx context.Context, req *UpdateTicketRequest) (*UpdateTicketResponse, error) {
	resp := new(UpdateTicketResponse)

	tx, err := h.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := h.Repo.WithTx(tx)

	// Do the nullable status
	var status repository.NullTicketStatus
	if req.Body.Status != nil {
		status.Valid = true
		status.TicketStatus = *req.Body.Status
	}

	// Do the basic update
	updateResult, err := qtx.UpdateTicketSimple(ctx, repository.UpdateTicketSimpleParams{
		ID:           req.TicketID,
		Status:       status,
		DepartmentID: req.Body.DepartmentID,
	})
	if err != nil {
		return nil, err
	}
	resp.Body = updateResult

	// Add new tags
	_, err = qtx.AddTagsToTicket(ctx, repository.AddTagsToTicketParams{
		Ticket: req.TicketID,
		Tags:   req.Body.AddTags,
	})
	if err != nil {
		return nil, err
	}

	// Remove some tags
	_, err = qtx.DeleteTagsFromTicket(ctx, repository.DeleteTagsFromTicketParams{
		Ticket: req.TicketID,
		Tags:   req.Body.RemoveTags,
	})
	if err != nil {
		return nil, err
	}

	tx.Commit(ctx)

	return resp, nil
}

// ==================== DELETE / HIDE ====================
type DeleteTicketRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeleteTicketResponse struct {
	Body repository.Ticket
}

func (h *TicketHandler) Delete(ctx context.Context, req *DeleteTicketRequest) (*DeleteTicketResponse, error) {
	resp := new(DeleteTicketResponse)

	ticket, err := h.Repo.DeleteTicket(ctx, repository.DeleteTicketParams{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}

	resp.Body = ticket
	return resp, nil
}

// =================== MERGE ====================
type MergeRequest struct {
	Body struct {
		Original   uuid.UUID   `json:"original"`
		Duplicates []uuid.UUID `json:"duplicates"`
	}
}
type MergeResonse struct {
}

func (h *TicketHandler) Merge(ctx context.Context, req *MergeRequest) (*MergeResonse, error) {
	resp := new(MergeResonse)

	tx, err := h.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := h.Repo.WithTx(tx)

	// Update the 'ticket' field for all details for a duplicates
	_, err = qtx.UpdateDetailsParent(ctx, repository.UpdateDetailsParentParams{
		Original:   req.Body.Original,
		Duplicates: req.Body.Duplicates,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to update details ticket ownership: %w", err)
	}

	// Set the average embedding
	err = qtx.MergeEmbeddings(ctx, repository.MergeEmbeddingsParams{
		Original:  req.Body.Original,
		Duplcates: req.Body.Duplicates,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to merge embeddings: %w", err)
	}

	// Delete obsolete tickets
	err = qtx.DeleteAfterMerge(ctx, repository.DeleteAfterMergeParams{
		Duplicates: req.Body.Duplicates,
	})
	if err != nil {
		return nil, err
	}

	tx.Commit(ctx)

	return resp, nil
}
