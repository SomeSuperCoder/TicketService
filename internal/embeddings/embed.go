package embeddings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pgvector/pgvector-go"
)

func GetEmbedding(text string) (*pgvector.Vector, error) {
	// Get the embedding via ollama
	emb, err := getEmbedding(text)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate the embedding due to: %w", err)
	}
	// Create a vector from the embedding
	vector := pgvector.NewVector(emb)

	return &vector, err
}

func getEmbedding(text string) ([]float32, error) {
	requestBody, _ := json.Marshal(map[string]string{
		"model":  "nomic-embed-text-v2-moe", // Russian-optimized model
		"prompt": text,
	})

	resp, err := http.Post("http://localhost:11434/api/embeddings",
		"application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Embedding []float32 `json:"embedding"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Embedding, nil
}
