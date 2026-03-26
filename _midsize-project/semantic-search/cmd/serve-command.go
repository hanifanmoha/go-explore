package cmd

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hanifanmoha/go-explore/_midsize-project/semantic-search/models"
	"github.com/hanifanmoha/go-explore/_midsize-project/semantic-search/pkg/helper"
	"github.com/hanifanmoha/go-explore/_midsize-project/semantic-search/services"
	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the semantic search server",
	Long:  `Start the semantic search server to handle search queries`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting the server...")
		startServer()
	},
}

func startServer() {
	mux := http.NewServeMux()

	handler, err := NewRouteHandler()
	if err != nil {
		log.Fatalf("Failed to initialize route handler: %v", err)
	}
	registerRouters(mux, handler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	errCh := make(chan error, 1)
	go func() {
		fmt.Printf("Server is running on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("server error: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	select {
	case sig := <-quit:
		fmt.Printf("Received signal: %v. Shutting down server...\n", sig)
	case err := <-errCh:
		fmt.Printf("Server error: %v\n", err)
	}

	handler.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	} else {
		fmt.Println("Server gracefully stopped")
	}
}

func registerRouters(mux *http.ServeMux, handler *RouteHandler) {
	mux.HandleFunc("/", handler.GetSearchPage)
	mux.HandleFunc("/{topic}", handler.GetTopicPage)
}

type RouteHandler struct {
	genaiSvc     *services.GenAIService
	dbConnection *services.DatabaseService
}

func (h *RouteHandler) Close() {
	if h.dbConnection != nil {
		h.dbConnection.Close()
	}
}

func NewRouteHandler() (*RouteHandler, error) {

	helper.LoadEnv()

	genaiSvc, err := services.NewGenAIService()
	if err != nil {
		log.Println("Error initializing GenAI service:", err)
		return nil, err
	}

	dbConnection, err := services.NewDatabaseService()
	if err != nil {
		log.Println("Error connecting to database:", err)
		return nil, err
	}

	return &RouteHandler{
		genaiSvc:     genaiSvc,
		dbConnection: dbConnection,
	}, nil
}

func (h *RouteHandler) GetSearchPage(w http.ResponseWriter, r *http.Request) {

	searchQuery := r.URL.Query().Get("q")

	log.Println("Query:", searchQuery)

	topicList, err := h.dbConnection.GetTopics(r.Context())
	if err != nil {
		log.Printf("Error fetching topic list: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var embeddings []models.Embedding
	if searchQuery != "" {
		vector, err := h.genaiSvc.GenerateEmbedding(searchQuery, services.TaskTypeRetrievalQuery)
		if err != nil {
			log.Printf("Error generating embedding for search query: %v\n", err)
			http.Error(w, "Failed to generate embedding for search query", http.StatusInternalServerError)
			return
		}

		embeddings, err = h.dbConnection.SearchEmbeddings(r.Context(), vector, 3)
		if err != nil {
			log.Printf("Error performing search: %v", err)
			http.Error(w, "Failed to perform search", http.StatusInternalServerError)
			return
		}
	}

	data := map[string]any{
		"SearchQuery": searchQuery,
		"TopicList":   topicList,
		"ItemList":    embeddings,
	}

	tmpl := template.Must(template.ParseFiles("views/base.html", "views/index.html"))

	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *RouteHandler) GetTopicPage(w http.ResponseWriter, r *http.Request) {

	topicName := r.PathValue("topic")

	topicList, err := h.dbConnection.GetTopics(r.Context())
	if err != nil {
		log.Printf("Error fetching topic list: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	embeddings, err := h.dbConnection.GetEmbeddingsByTopic(r.Context(), topicName)
	if err != nil {
		log.Printf("Error fetching topic data: %v", err)
		http.Error(w, "Failed to fetch topic data", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("views/base.html", "views/topic.html"))

	data := map[string]any{
		"TopicList": topicList,
		"ItemList":  embeddings,
	}

	if err := tmpl.ExecuteTemplate(w, "topic.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
