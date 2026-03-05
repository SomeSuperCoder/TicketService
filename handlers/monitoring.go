package handlers

import (
	"context"
	"fmt"

	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/google/uuid"
)

type MonitoringHandler struct {
	Repo *repository.Queries
}

// ==================== REQUEST/RESPONSE TYPES ====================

type GetOverdueRequest struct {
	DepartmentID int32 `query:"department_id"`
	MinLostDays  int32 `query:"min_lost_days" default:"1"`
	Limit        int32 `query:"limit" default:"50"`
}

type OverdueTicket struct {
	ID              uuid.UUID `json:"id"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	SubcategoryID   int32     `json:"subcategory_id"`
	DepartmentID    *int32    `json:"department_id"`
	StatusStartDate string    `json:"status_start_date"`
	LostDays        int32     `json:"lost_days"`
}

type GetOverdueResponse struct {
	Body struct {
		Tickets []OverdueTicket `json:"tickets"`
		Total   int64           `json:"total"`
	}
}

// ==================== HANDLER METHODS ====================

func (h *MonitoringHandler) GetOverdue(ctx context.Context, req *GetOverdueRequest) (*GetOverdueResponse, error) {
	resp := new(GetOverdueResponse)

	// Validate query parameters
	if err := validateOverdueParams(req); err != nil {
		return nil, err
	}

	// Convert DepartmentID to pointer (0 means not provided)
	var departmentID *int32
	if req.DepartmentID > 0 {
		departmentID = &req.DepartmentID
	}

	// Fetch overdue tickets from repository
	tickets, err := h.Repo.GetOverdueTickets(ctx, repository.GetOverdueTicketsParams{
		MinLostDays:  req.MinLostDays,
		DepartmentID: departmentID,
		Limit:        req.Limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch overdue tickets: %w", err)
	}

	// Fetch total count of overdue tickets
	total, err := h.Repo.CountOverdueTickets(ctx, repository.CountOverdueTicketsParams{
		MinLostDays:  req.MinLostDays,
		DepartmentID: departmentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to count overdue tickets: %w", err)
	}

	// Format response
	resp.Body.Tickets = formatTickets(tickets)
	resp.Body.Total = total

	return resp, nil
}

// ==================== VALIDATION ====================

func validateOverdueParams(req *GetOverdueRequest) error {
	// Validate min_lost_days is not negative
	if req.MinLostDays < 0 {
		return fmt.Errorf("min_lost_days must be a positive integer")
	}

	// Validate limit is not negative and does not exceed maximum
	if req.Limit < 0 {
		return fmt.Errorf("limit must be a positive integer")
	}
	if req.Limit > 100 {
		return fmt.Errorf("limit cannot exceed 100")
	}

	// Validate department_id if provided (0 means not provided)
	if req.DepartmentID < 0 {
		return fmt.Errorf("department_id must be a positive integer")
	}

	return nil
}

// ==================== FORMATTING ====================

func formatTickets(tickets []repository.GetOverdueTicketsRow) []OverdueTicket {
	result := make([]OverdueTicket, len(tickets))
	for i, ticket := range tickets {
		result[i] = OverdueTicket{
			ID:              ticket.ID,
			Description:     ticket.Description,
			Status:          string(ticket.Status),
			SubcategoryID:   ticket.SubcategoryID,
			DepartmentID:    ticket.DepartmentID,
			StatusStartDate: ticket.StatusStartDate.Format("2006-01-02T15:04:05Z07:00"),
			LostDays:        ticket.LostDays,
		}
	}
	return result
}
