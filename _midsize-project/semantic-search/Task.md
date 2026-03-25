# General Description

- This is a project to discover implement semantic search with embedding
- I will use Gemini API to do embedding
- I have prepared a dataset of animal descriptions in YAML format, in /datasource/animals.yaml
- This project will be implemented in golang

# Tech Stack

- Golang as service to serve UI & API
  - Go HTML for UI
  - Go Gin for API
  - Cobra for CLI
- Gemini API for embedding
- YAML for data source
- PostgreSQL for storing data and embeddings

# Implementations

## Embedding Steps

- Use data in /datasource/animals.yaml as datasource
- Create a CLI with cobra to create embedding and store in PostgreSQL
  - Also, save the vector data in the original yaml file in the animal records

## API Steps : Documentation UI

- Create a UI with several tabs in the top. for each tab is a topic from datasource.
- URL: /<topic / filename of datasource yml>#<item title>
- When the user click each tab, it will redirect to tha page containing all data in the topic. And if provided with anchor (#item_name), it will scroll to the item in the page.
- Make the UI simple, just make it centered in the small screen, and make it look like a documentation page. No need to make it fancy.

## API Steps : Search API

- Create an API endpoint to search for relevant items based on a query.
- In the top tabs list, add a search tab in the most left.
- URL: /search?q=<query>
- When the user click the search tab, it will redirect to the search page. The search page will display the search results based on the query. The search results will be sorted by relevance, and each result will display the item title, a snippet of the description, and a link to the item page.
- If the item click, it will be redirect to /<topic / filename of datasource yml>#<item title>
- if the URL is not provided with query, only /search, it will display the search page with empty result, with an input text area and a search button. When the user input a query and click the search button, it will redirect to /search?q=<query> and display the search results.

# Implementation Steps

# Step 1: Set Up Project Structure and Dependencies

- Initialize the Go module (already done with `go mod init github.com/hanifanmoha/go-explore/_midsize-project/semantic-search`).
- Install required dependencies:
  - `go get github.com/spf13/cobra@latest` for CLI.
  - `go get github.com/gin-gonic/gin` for API.
  - `go get gopkg.in/yaml.v3` for YAML parsing.
  - `go get github.com/lib/pq` for PostgreSQL driver.
  - `go get github.com/google/generative-ai-go/genai` for Gemini API (assuming the package name).
- Set up the project directory structure:
  - `cmd/` for CLI commands.
  - `internal/` for internal packages (e.g., `embedding/`, `database/`, `api/`).
  - `web/` for HTML templates.
  - `datasource/` (already exists).
- Create a `config.yaml` or environment variables for API keys and database connection.

# Step 2: Implement Database Schema and Connection

- Create a `docker-compose.yml` file to set up PostgreSQL with the pgvector extension for vector storage:
  ```yaml
  version: '3.8'
  services:
    postgres:
      image: pgvector/pgvector:pg16
      environment:
        POSTGRES_DB: semantic_search
        POSTGRES_USER: user
        POSTGRES_PASSWORD: password
      ports:
        - "5432:5432"
      volumes:
        - postgres_data:/var/lib/postgresql/data
  volumes:
    postgres_data:
  ```
- Design the PostgreSQL schema:
  - Table `items` with columns: `id` (serial primary key), `topic` (text), `title` (text), `description` (text), `embedding` (vector(768) for Gemini embeddings).
- Create a database package in `internal/database/` to handle connections and migrations.
- Use SQL scripts or Go code to create tables.
- Implement functions to insert and query items with embeddings.

# Step 3: Implement Embedding CLI with Cobra

- Create the main CLI entry point in `cmd/cli/main.go`.
- Use Cobra to define a command like `embed` that:
  - Reads the YAML files from `datasource/`.
  - For each item, calls Gemini API to generate embeddings.
  - Stores the embedding in PostgreSQL.
  - Updates the YAML file to include the embedding vector in each item record.
- Handle batch processing and error handling.
- Create a `Makefile` for building and running the CLI:
  ```
  .PHONY: build-cli run-cli clean

  build-cli:
  	go build -o bin/cli cmd/cli/main.go

  run-cli:
  	./bin/cli

  clean:
  	rm -rf bin/
  ```
- Example pseudo code:
  ```go
  func embedCommand() {
      files := listYamlFiles("datasource/")
      for _, file := range files {
          data := parseYaml(file)
          for _, item := range data.Items {
              vector := callGeminiEmbedding(item.Description)
              saveToDB(item, vector)
              updateYaml(file, item, vector)
          }
      }
  }
  ```

# Step 4: Implement the API Server with Gin

- Create the main server in `cmd/server/main.go`.
- Set up Gin router.
- Implement routes:
  - `/<topic>`: Serve documentation UI for the topic.
  - `/<topic>#<item>`: Same as above, with anchor.
  - `/search`: Handle search page and results.
- Use HTML templates in `web/templates/` for rendering pages.

# Step 5: Implement Documentation UI

- Create HTML templates for the UI:
  - `layout.html`: Base template with tabs for topics.
  - `topic.html`: Display all items in a topic, with anchors.
- In the API handler for `/<topic>`, parse the YAML, render the template with data.
- Make the UI simple: centered, documentation-style, with tabs at the top.
- For search tab, link to `/search`.

# Step 6: Implement Search API and UI

- Create a search handler for `/search`:
  - If no query, render search form.
  - If query, generate embedding for the query using Gemini.
  - Query PostgreSQL for similar items using vector similarity (e.g., cosine distance).
  - Sort results by relevance.
- Render search results page with titles, snippets, and links to item pages.
- Use templates like `search.html` for the search interface and results.

# Step 7: Testing and Validation

- Test the CLI embedding process with sample data.
- Test API endpoints with curl or browser.
- Validate embeddings by checking search relevance.
- Ensure YAML updates correctly with vectors.

# Step 8: Deployment and Final Touches

- Add configuration for production (e.g., environment variables for API keys).
- Implement logging and error handling.
- Document the project in README.md.
- Run the full pipeline: embed data, start server, test search.