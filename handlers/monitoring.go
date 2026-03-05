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

type DepartmentEfficiency struct {
	DepartmentName    string  `json:"department_name"`
	InProgress        int64   `json:"in_progress"`
	Overdue           int64   `json:"overdue"`
	AvgResolutionDays float64 `json:"avg_resolution_days"`
	TrendPercent      float64 `json:"trend_percent"`
}

type GetDepartmentEfficiencyResponse struct {
	Body struct {
		Departments []DepartmentEfficiency `json:"departments"`
	}
}

type GetDepartmentEfficiencyRequest struct {
}

type KPIMetrics struct {
	AvgResponseDays   float64 `json:"avg_response_days"`
	OverdueCount      int64   `json:"overdue_count"`
	SatisfactionIndex float64 `json:"satisfaction_index"`
}

type GetKPIRequest struct {
	Period string `query:"period" enum:"week,month,year" default:"month"`
}

type GetKPIResponse struct {
	Body KPIMetrics
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

func (h *MonitoringHandler) GetDepartmentEfficiency(ctx context.Context, req *GetDepartmentEfficiencyRequest) (*GetDepartmentEfficiencyResponse, error) {
	resp := new(GetDepartmentEfficiencyResponse)

	// Fetch department efficiency data from repository
	departments, err := h.Repo.GetDepartmentEfficiency(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch department efficiency: %w", err)
	}

	// Format response
	resp.Body.Departments = formatDepartmentEfficiency(departments)

	return resp, nil
}

func (h *MonitoringHandler) GetKPI(ctx context.Context, req *GetKPIRequest) (*GetKPIResponse, error) {
	resp := new(GetKPIResponse)

	// Validate period parameter
	if err := validateKPIPeriod(req); err != nil {
		return nil, err
	}

	// Fetch KPI data from repository
	kpi, err := h.Repo.GetKPI(ctx, repository.GetKPIParams{
		Period: req.Period,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch KPI metrics: %w", err)
	}

	// Format response
	satisfactionIndex := 0.0
	if kpi.SatisfactionIndex != nil {
		if val, ok := kpi.SatisfactionIndex.(float64); ok {
			satisfactionIndex = val
		}
	}

	resp.Body = KPIMetrics{
		AvgResponseDays:   kpi.AvgResponseDays,
		OverdueCount:      kpi.OverdueCount,
		SatisfactionIndex: satisfactionIndex,
	}

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

func validateKPIPeriod(req *GetKPIRequest) error {
	validPeriods := map[string]bool{
		"week":  true,
		"month": true,
		"year":  true,
	}

	if !validPeriods[req.Period] {
		return fmt.Errorf("period must be one of: week, month, year")
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

func formatDepartmentEfficiency(departments []repository.GetDepartmentEfficiencyRow) []DepartmentEfficiency {
	result := make([]DepartmentEfficiency, len(departments))
	for i, dept := range departments {
		trendPercent := 0.0
		if dept.TrendPercent != nil {
			if val, ok := dept.TrendPercent.(float64); ok {
				trendPercent = val
			}
		}

		result[i] = DepartmentEfficiency{
			DepartmentName:    dept.DepartmentName,
			InProgress:        dept.InProgress,
			Overdue:           dept.Overdue,
			AvgResolutionDays: dept.AvgResolutionDays,
			TrendPercent:      trendPercent,
		}
	}
	return result
}
