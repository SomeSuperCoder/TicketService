package embeddings

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func GetEmbedding(text string) ([]float32, error) {
	requestBody, _ := json.Marshal(map[string]string{
		"model":  "evilfreelancer/enbeddrus", // Russian-optimized model
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
