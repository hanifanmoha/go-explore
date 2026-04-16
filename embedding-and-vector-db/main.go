package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/genai"
)

const (
	API_KEY = "<API_KEY_FROM_https://aistudio.google.com/u/0/api-keys>"
)

func main() {
	ctx := context.Background()

	cfg := genai.ClientConfig{
		APIKey: API_KEY,
	}
	client, err := genai.NewClient(ctx, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	text := "The earth is flat."
	outputDimension := int32(100)

	result, err := client.Models.EmbedContent(
		ctx,
		"gemini-embedding-001",
		[]*genai.Content{genai.NewContentFromText(text, "")},
		&genai.EmbedContentConfig{
			OutputDimensionality: &outputDimension,
			TaskType:             "SEMANTIC_SIMILARITY",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	if len(result.Embeddings) == 0 {
		log.Fatal("no embeddings returned")
	}

	vectorValue := result.Embeddings[0].Values
	fmt.Println("Result (first 10 values): ", vectorValue[:10]) // Print the first 10 values of the embedding vector
	fmt.Println("Length: ", len(vectorValue))

}
