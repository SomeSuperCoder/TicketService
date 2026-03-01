package handlers

import (
	"context"
	"fmt"

	"github.com/SomeSuperCoder/OnlineShop/internal/embeddings"
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

	// Get the embedding via ollama
	emb, err := embeddings.GetEmbedding("The quick brown fox jumps over the lazy dog")
	if err != nil {
		return nil, fmt.Errorf("Failed to generate the embedding due to: %w", err)
	}
	// Create a vector from the embedding
	vector := pgvector.NewVector(emb)

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
