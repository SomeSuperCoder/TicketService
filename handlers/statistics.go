package handlers

import (
	"context"
	"time"

	"github.com/SomeSuperCoder/OnlineShop/repository"
)

type StatisticsHandler struct {
	Repo *repository.Queries
}

type GetSummaryResponse struct {
	Body repository.GetStatisticsSummaryRow
}

func (h *StatisticsHandler) GetSummary(ctx context.Context, req *struct{}) (*GetSummaryResponse, error) {
	resp := new(GetSummaryResponse)

	result, err := h.Repo.GetStatisticsSummary(ctx)
	if err != nil {
		return nil, err
	}

	resp.Body = result

	return resp, nil
}

// ==========================================

type GetCategoryStatisticsResponse struct {
	Body []repository.GetCategoryStatisticsRow
}

func (h *StatisticsHandler) GetCategoryStatistics(ctx context.Context, req *struct{}) (*GetCategoryStatisticsResponse, error) {
	resp := new(GetCategoryStatisticsResponse)

	result, err := h.Repo.GetCategoryStatistics(ctx)
	if err != nil {
		return nil, err
	}

	resp.Body = result
	return resp, nil
}

// ==========================================

type GetDynamicsRequest struct {
	Period    string `query:"period" enum:"day,week" default:"day"`
	StartDate string `query:"start_date"` // Format: YYYY-MM-DD
	EndDate   string `query:"end_date"`   // Format: YYYY-MM-DD
}

type DynamicsDataPoint struct {
	Date     string `json:"date"`
	Received int64  `json:"received"`
	Resolved int64  `json:"resolved"`
}

type GetDynamicsResponse struct {
	Body struct {
		Period string              `json:"period"`
		Data   []DynamicsDataPoint `json:"data"`
	}
}

func (h *StatisticsHandler) GetDynamics(ctx context.Context, req *GetDynamicsRequest) (*GetDynamicsResponse, error) {
	resp := new(GetDynamicsResponse)
	resp.Body.Period = req.Period

	// Parse dates with defaults
	var startDate, endDate time.Time
	var err error

	if req.StartDate == "" {
		// Default to 30 days ago
		startDate = time.Now().AddDate(0, 0, -30)
	} else {
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, err
		}
	}

	if req.EndDate == "" {
		// Default to now
		endDate = time.Now()
	} else {
		endDate, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, err
		}
	}

	if req.Period == "week" {
		// Get weekly dynamics
		result, err := h.Repo.GetTicketDynamicsByWeek(ctx, repository.GetTicketDynamicsByWeekParams{
			StartDate: startDate,
			EndDate:   endDate,
		})
		if err != nil {
			return nil, err
		}

		resp.Body.Data = make([]DynamicsDataPoint, len(result))
		for i, row := range result {
			dateStr := ""
			if row.WeekStart.Valid {
				dateStr = row.WeekStart.Time.Format("2006-01-02")
			}
			resp.Body.Data[i] = DynamicsDataPoint{
				Date:     dateStr,
				Received: row.Received,
				Resolved: row.Resolved,
			}
		}
	} else {
		// Get daily dynamics (default)
		result, err := h.Repo.GetTicketDynamics(ctx, repository.GetTicketDynamicsParams{
			StartDate: startDate,
			EndDate:   endDate,
		})
		if err != nil {
			return nil, err
		}

		resp.Body.Data = make([]DynamicsDataPoint, len(result))
		for i, row := range result {
			dateStr := ""
			if row.Date.Valid {
				dateStr = row.Date.Time.Format("2006-01-02")
			}
			resp.Body.Data[i] = DynamicsDataPoint{
				Date:     dateStr,
				Received: row.Received,
				Resolved: row.Resolved,
			}
		}
	}

	return resp, nil
}
