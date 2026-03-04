package handlers

import (
	"context"

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
