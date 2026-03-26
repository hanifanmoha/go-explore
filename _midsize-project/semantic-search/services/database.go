package services

import (
	"context"
	"fmt"
	"log"

	"github.com/hanifanmoha/go-explore/_midsize-project/semantic-search/models"
	"github.com/hanifanmoha/go-explore/_midsize-project/semantic-search/pkg/helper"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseService struct {
	Pool *pgxpool.Pool
}

func NewDatabaseService() (*DatabaseService, error) {

	dbURL := helper.GetDBURL()

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	dbSvc := &DatabaseService{Pool: pool}

	// Run migrations on startup
	if err := dbSvc.initializeSchema(); err != nil {
		log.Printf("Warning: failed to initialize schema: %v", err)
	}

	return dbSvc, nil
}

func (d *DatabaseService) Close() {
	d.Pool.Close()
}

func (d *DatabaseService) InsertEmbedding(topic, name, description string, embedding []float32) error {

	ctx := context.Background()

	embeddingJSON, err := helper.EmbeddingToJSON(embedding)
	if err != nil {
		return fmt.Errorf("failed to convert embedding to JSON: %w", err)
	}

	_, err = d.Pool.Exec(ctx, `
				INSERT INTO embeddings (topic, name, description, embedding)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (topic, name) DO UPDATE SET
				description = EXCLUDED.description,
				embedding = EXCLUDED.embedding,
				updated_at = CURRENT_TIMESTAMP
	`, topic, name, description, embeddingJSON)
	if err != nil {
		return fmt.Errorf("failed to insert embedding: %w", err)
	}

	return nil
}

func (d *DatabaseService) SearchEmbeddings(ctx context.Context, queryEmbedding []float32, limit int) ([]models.Embedding, error) {

	embeddingJSON, err := helper.EmbeddingToJSON(queryEmbedding)
	if err != nil {
		return nil, fmt.Errorf("failed to convert query embedding to JSON: %w", err)
	}

	rows, err := d.Pool.Query(ctx, `
		SELECT id, topic, name, description, embedding
		FROM embeddings
		ORDER BY embedding <=> $1::vector
		LIMIT $2
	`, embeddingJSON, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w", err)
	}
	defer rows.Close()

	var results []models.Embedding
	for rows.Next() {
		var e models.Embedding
		var embeddingJSON string

		if err := rows.Scan(&e.ID, &e.Topic, &e.Name, &e.Description, &embeddingJSON); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		e.Embedding, err = helper.JSONToEmbedding(embeddingJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to convert embedding JSON to slice: %w", err)
		}

		results = append(results, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return results, nil
}

func (d *DatabaseService) GetTopics(ctx context.Context) ([]string, error) {
	rows, err := d.Pool.Query(ctx, `
		SELECT DISTINCT topic
		FROM embeddings
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get topics query: %w", err)
	}
	defer rows.Close()

	var topics []string
	for rows.Next() {
		var topic string
		if err := rows.Scan(&topic); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		topics = append(topics, topic)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return topics, nil
}

func (d *DatabaseService) GetEmbeddingsByTopic(ctx context.Context, topic string) ([]models.Embedding, error) {
	rows, err := d.Pool.Query(ctx, `
		SELECT id, topic, name, description, embedding
		FROM embeddings
		WHERE topic = $1
	`, topic)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get embeddings by topic query: %w", err)
	}
	defer rows.Close()

	var results []models.Embedding
	for rows.Next() {
		var e models.Embedding
		var embeddingJSON string

		if err := rows.Scan(&e.ID, &e.Topic, &e.Name, &e.Description, &embeddingJSON); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		e.Embedding, err = helper.JSONToEmbedding(embeddingJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to convert embedding JSON to slice: %w", err)
		}

		results = append(results, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return results, nil
}

func (d *DatabaseService) initializeSchema() error {
	ctx := context.Background()

	// Enable pgvector extension
	if _, err := d.Pool.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector"); err != nil {
		log.Printf("Warning: failed to create vector extension: %v", err)
	}

	// Create embeddings table
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS embeddings (
            id SERIAL PRIMARY KEY,
            topic VARCHAR(255) NOT NULL,
            name VARCHAR(255) NOT NULL,
            description TEXT NOT NULL,
            embedding vector(1536) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            UNIQUE(topic, name)
        )
    `

	if _, err := d.Pool.Exec(ctx, createTableQuery); err != nil {
		log.Printf("Warning: failed to create embeddings table: %v", err)
	}

	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_embeddings_topic ON embeddings(topic)",
		"CREATE INDEX IF NOT EXISTS idx_embeddings_name ON embeddings(name)",
		"CREATE INDEX IF NOT EXISTS idx_embeddings_vector ON embeddings USING ivfflat (embedding vector_cosine_ops)",
	}

	for _, indexQuery := range indexes {
		if _, err := d.Pool.Exec(ctx, indexQuery); err != nil {
			log.Printf("Warning: failed to create index: %v", err)
		}
	}

	log.Println("Database schema initialized successfully")
	return nil
}
