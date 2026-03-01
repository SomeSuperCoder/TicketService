package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/pgvector/pgvector-go"
)

type TicketHandler struct {
	Repo *repository.Queries
}

type PostTicketRequest struct{}
type PostTicketResponse struct {
	Body repository.Ticket
}

func (h *TicketHandler) Post(ctx context.Context, req *PostTicketRequest) (*PostTicketResponse, error) {
	resp := new(PostTicketResponse)

	vector := pgvector.NewVector(make([]float32, 1536))

	result, err := h.Repo.CreateTicket(ctx, repository.CreateTicketParams{
		Title:       "This is a title",
		Description: "This is a description",
		Longitude:   33,
		Latitude:    44,
		Embedding:   &vector,
	})
	resp.Body = result

	return resp, err
}
