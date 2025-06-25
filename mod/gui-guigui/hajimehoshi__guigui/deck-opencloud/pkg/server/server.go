package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"opencloud/pkg/indexer"
	"opencloud/pkg/search"
)

// Config holds server configuration
type Config struct {
	Host     string
	Port     string
	IndexDir string
	DataDir  string
	Debug    bool
}

// Server represents the OpenCloud collaboration server
type Server struct {
	config   *Config
	indexer  *indexer.Indexer
	searcher *search.Searcher
	server   *http.Server
}

// SearchRequest represents an API search request
type SearchRequest struct {
	Query   string            `json:"query"`
	Limit   int               `json:"limit,omitempty"`
	Offset  int               `json:"offset,omitempty"`
	Filters map[string]string `json:"filters,omitempty"`
}

// SearchResponse represents an API search response
type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
	Took    int            `json:"took"` // milliseconds
}

// SearchResult represents a search result
type SearchResult struct {
	ID       string                 `json:"id"`
	Score    float64                `json:"score"`
	Title    string                 `json:"title"`
	Path     string                 `json:"path"`
	Type     string                 `json:"type"`
	Snippet  string                 `json:"snippet"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// IndexRequest represents a document indexing request
type IndexRequest struct {
	Path     string                 `json:"path"`
	Content  string                 `json:"content,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// IndexResponse represents an indexing response
type IndexResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// New creates a new OpenCloud server
func New(config *Config) (*Server, error) {
	srv := &Server{
		config:   config,
		indexer:  indexer.New(config.IndexDir),
		searcher: search.New(config.IndexDir),
	}

	// Create HTTP server
	mux := http.NewServeMux()
	srv.setupRoutes(mux)

	srv.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
		Handler:      srv.corsMiddleware(srv.loggingMiddleware(mux)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return srv, nil
}

// setupRoutes configures HTTP routes
func (s *Server) setupRoutes(mux *http.ServeMux) {
	// Health check
	mux.HandleFunc("/health", s.handleHealth)

	// API endpoints
	mux.HandleFunc("/api/search", s.handleSearch)
	mux.HandleFunc("/api/index", s.handleIndex)
	mux.HandleFunc("/api/documents", s.handleDocuments)

	// Web interface
	mux.HandleFunc("/", s.handleWebInterface)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
		"services": map[string]string{
			"search":  "available",
			"indexer": "available",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSearch handles search requests
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SearchRequest

	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else {
		// GET request - parse query parameters
		req.Query = r.URL.Query().Get("q")
		if req.Query == "" {
			req.Query = r.URL.Query().Get("query")
		}

		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if limit, err := strconv.Atoi(limitStr); err == nil {
				req.Limit = limit
			}
		}
	}

	if req.Query == "" {
		http.Error(w, "Query parameter required", http.StatusBadRequest)
		return
	}

	if req.Limit <= 0 {
		req.Limit = 50 // default
	}

	start := time.Now()
	results, err := s.searcher.Search(req.Query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Search failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert results
	apiResults := make([]SearchResult, 0, len(results))
	for _, result := range results {
		apiResults = append(apiResults, SearchResult{
			ID:       result.ID,
			Score:    result.Score,
			Title:    result.Title,
			Path:     result.Path,
			Type:     result.Type,
			Snippet:  result.Snippet,
			Metadata: result.Metadata,
		})
	}

	response := SearchResponse{
		Results: apiResults,
		Total:   len(apiResults),
		Took:    int(time.Since(start).Milliseconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleIndex handles document indexing requests
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req IndexRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Path == "" {
		http.Error(w, "Path required", http.StatusBadRequest)
		return
	}

	// For now, just trigger re-indexing of the data directory
	// In a real implementation, this would index the specific document
	err := s.indexer.IndexDirectory(s.config.DataDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Indexing failed: %v", err), http.StatusInternalServerError)
		return
	}

	response := IndexResponse{
		ID:     filepath.Base(req.Path),
		Status: "indexed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleDocuments handles document listing requests
func (s *Server) handleDocuments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simple implementation - search for all documents
	results, err := s.searcher.Search("*")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list documents: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"documents": results,
		"total":     len(results),
	})
}

// handleWebInterface serves a simple web interface
func (s *Server) handleWebInterface(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	html := `<!DOCTYPE html>
<html>
<head>
    <title>OpenCloud - Collaboration Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .header { color: #2c3e50; }
        .search-box { margin: 20px 0; }
        .search-box input { padding: 10px; width: 300px; }
        .search-box button { padding: 10px 20px; }
        .results { margin-top: 20px; }
        .result { border: 1px solid #ddd; margin: 10px 0; padding: 15px; }
        .result-title { font-weight: bold; color: #2980b9; }
        .result-path { color: #7f8c8d; font-size: 0.9em; }
        .result-snippet { margin-top: 5px; }
    </style>
</head>
<body>
    <h1 class="header">🌩️ OpenCloud - Collaboration Server</h1>
    <p>Search and collaboration server with Markdown and DeckSh support</p>
    
    <div class="search-box">
        <input type="text" id="searchQuery" placeholder="Enter search query..." />
        <button onclick="search()">Search</button>
    </div>
    
    <div id="results" class="results"></div>
    
    <script>
        function search() {
            const query = document.getElementById('searchQuery').value;
            if (!query) return;
            
            fetch('/api/search?q=' + encodeURIComponent(query))
                .then(response => response.json())
                .then(data => {
                    const resultsDiv = document.getElementById('results');
                    if (data.results && data.results.length > 0) {
                        resultsDiv.innerHTML = '<h3>Found ' + data.total + ' results (' + data.took + 'ms)</h3>' +
                            data.results.map(result => 
                                '<div class="result">' +
                                '<div class="result-title">' + result.title + '</div>' +
                                '<div class="result-path">' + result.path + ' (' + result.type + ') - Score: ' + result.score.toFixed(3) + '</div>' +
                                '<div class="result-snippet">' + (result.snippet || '') + '</div>' +
                                '</div>'
                            ).join('');
                    } else {
                        resultsDiv.innerHTML = '<p>No results found</p>';
                    }
                })
                .catch(error => {
                    document.getElementById('results').innerHTML = '<p>Error: ' + error + '</p>';
                });
        }
        
        document.getElementById('searchQuery').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') search();
        });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs HTTP requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		if s.config.Debug {
			log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
		}
	})
}
