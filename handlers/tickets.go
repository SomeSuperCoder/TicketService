package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/internal/embeddings"
	"github.com/SomeSuperCoder/OnlineShop/repository"
)

type TicketHandler struct {
	Repo *repository.Queries
}

type PostTicketRequest struct {
	Body struct {
		MockText string `json:"mock_text"`
	}
}
type PostTicketResponse struct {
	Body repository.Ticket
}

func (h *TicketHandler) Post(ctx context.Context, req *PostTicketRequest) (*PostTicketResponse, error) {
	resp := new(PostTicketResponse)

	vector, err := embeddings.GetEmbedding(req.Body.MockText)
	if err != nil {
		return nil, err
	}

	result, err := h.Repo.CreateTicket(ctx, repository.CreateTicketParams{
		Title:       req.Body.MockText,
		Description: req.Body.MockText,
		Longitude:   33,
		Latitude:    44,
		Embedding:   vector,
	})
	resp.Body = result

	return resp, err
}

type SearchByMeaningRequest struct {
	Query string `query:"query"`
}
type SearchByMeaningResponse struct {
	Body struct {
		Tickets []repository.SearchTicketsByMeaningRow `json:"tickets"`
	}
}

func (h *TicketHandler) SearchByMeaning(ctx context.Context, req *SearchByMeaningRequest) (*SearchByMeaningResponse, error) {
	resp := new(SearchByMeaningResponse)

	vector, err := embeddings.GetEmbedding(req.Query)
	if err != nil {
		return nil, err
	}

	result, err := h.Repo.SearchTicketsByMeaning(ctx, repository.SearchTicketsByMeaningParams{
		QueryEmbedding: vector,
	})
	resp.Body.Tickets = result

	return resp, nil
}
