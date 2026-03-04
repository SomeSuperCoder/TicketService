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
