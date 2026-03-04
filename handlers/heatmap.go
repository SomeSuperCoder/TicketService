package handlers

import (
	"context"
	"time"

	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HeatmapHandler struct {
	Repo *repository.Queries
	Pool *pgxpool.Pool
}

// ==================== GET HEATMAP POINTS ====================

type GetHeatmapPointsRequest struct {
	Period     string  `query:"period" enum:"week,month,year" default:"month"`
	Categories []int32 `query:"categories"`
}

type HeatmapPoint struct {
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Intensity int64   `json:"intensity"`
	Count     int64   `json:"count"`
}

type GetHeatmapPointsResponse struct {
	Body struct {
		Points []HeatmapPoint `json:"points"`
	}
}

func (h *HeatmapHandler) GetPoints(ctx context.Context, req *GetHeatmapPointsRequest) (*GetHeatmapPointsResponse, error) {
	resp := new(GetHeatmapPointsResponse)

	// Calculate date range based on period
	endDate := time.Now()
	var startDate time.Time

	switch req.Period {
	case "week":
		startDate = endDate.AddDate(0, 0, -7)
	case "year":
		startDate = endDate.AddDate(-1, 0, 0)
	default: // month
		startDate = endDate.AddDate(0, -1, 0)
	}

	// Get heatmap points
	points, err := h.Repo.GetHeatmapPoints(ctx, repository.GetHeatmapPointsParams{
		StartDate:  startDate,
		EndDate:    endDate,
		Categories: req.Categories,
	})
	if err != nil {
		return nil, err
	}

	// Convert to response format
	resp.Body.Points = make([]HeatmapPoint, len(points))
	for i, point := range points {
		resp.Body.Points[i] = HeatmapPoint{
			Lat:       point.Lat,
			Lng:       point.Lng,
			Intensity: point.Intensity,
			Count:     point.Count,
		}
	}

	return resp, nil
}

// ==================== GET HEATMAP STATS ====================

type ProblemLocation struct {
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	TicketCount int64   `json:"ticket_count"`
}

type GetHeatmapStatsResponse struct {
	Body struct {
		TopLocations          []ProblemLocation `json:"top_locations"`
		TotalLocations        int64             `json:"total_locations"`
		AvgTicketsPerLocation float64           `json:"avg_tickets_per_location"`
	}
}

func (h *HeatmapHandler) GetStats(ctx context.Context, req *struct{}) (*GetHeatmapStatsResponse, error) {
	resp := new(GetHeatmapStatsResponse)

	// Get top problem locations
	topLocations, err := h.Repo.GetTopProblemLocations(ctx)
	if err != nil {
		return nil, err
	}

	resp.Body.TopLocations = make([]ProblemLocation, len(topLocations))
	for i, loc := range topLocations {
		resp.Body.TopLocations[i] = ProblemLocation{
			Lat:         loc.Lat,
			Lng:         loc.Lng,
			TicketCount: loc.TicketCount,
		}
	}

	// Get overall stats
	stats, err := h.Repo.GetHeatmapStats(ctx)
	if err != nil {
		return nil, err
	}

	resp.Body.TotalLocations = stats.TotalLocations
	// AvgTicketsPerLocation is int32 from sqlc, convert to float64
	resp.Body.AvgTicketsPerLocation = float64(stats.AvgTicketsPerLocation)

	return resp, nil
}
