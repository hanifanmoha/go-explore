package models

type Embedding struct {
	ID          int       `json:"id"`
	Topic       string    `json:"topic"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Embedding   []float32 `json:"embedding"`
}
