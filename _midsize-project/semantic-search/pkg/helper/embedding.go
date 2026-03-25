package helper

import (
	"encoding/json"
	"fmt"
)

func EmbeddingToJSON(embedding []float32) (string, error) {

	embeddingJSON, err := json.Marshal(embedding)
	if err != nil {
		return "", fmt.Errorf("failed to marshal embedding: %w", err)
	}

	return string(embeddingJSON), nil
}

func JSONToEmbedding(embeddingJSON string) ([]float32, error) {
	var embedding []float32
	err := json.Unmarshal([]byte(embeddingJSON), &embedding)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal embedding: %w", err)
	}
	return embedding, nil
}
