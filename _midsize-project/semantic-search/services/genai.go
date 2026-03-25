package services

import (
	"context"
	"fmt"

	"github.com/hanifanmoha/go-explore/_midsize-project/semantic-search/pkg/helper"
	"google.golang.org/genai"
)

const EmbeddingModel = "gemini-embedding-001"
const TaskTypeRetrievalDocument = "RETRIEVAL_DOCUMENT"
const TaskTypeRetrievalQuery = "RETRIEVAL_QUERY"

type GenAIService struct {
	Client *genai.Client
}

func NewGenAIService() (*GenAIService, error) {
	ctx := context.Background()

	cfg := genai.ClientConfig{
		APIKey: helper.GetGenAIAPIKey(),
	}

	client, err := genai.NewClient(ctx, &cfg)
	if err != nil {
		return nil, err
	}

	return &GenAIService{
		Client: client,
	}, nil
}

func (s *GenAIService) GenerateEmbedding(text string, taskType string) ([]float32, error) {

	var OutputDimensionality int32 = 1536

	result, err := s.Client.Models.EmbedContent(
		context.Background(),
		EmbeddingModel,
		[]*genai.Content{genai.NewContentFromText(text, "")},
		&genai.EmbedContentConfig{
			OutputDimensionality: &OutputDimensionality,
			TaskType:             taskType,
		},
	)
	if err != nil {
		return nil, err
	}

	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return result.Embeddings[0].Values, nil
}
