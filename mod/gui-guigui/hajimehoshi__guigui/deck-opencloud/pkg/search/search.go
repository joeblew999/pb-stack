package search

import (
	"fmt"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

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

// Searcher handles search operations
type Searcher struct {
	index    bleve.Index
	indexDir string
}

// New creates a new searcher
func New(indexDir string) *Searcher {
	return &Searcher{
		indexDir: indexDir,
	}
}

// Search performs a search query
func (s *Searcher) Search(queryStr string) ([]*SearchResult, error) {
	// Open index
	index, err := bleve.Open(s.indexDir)
	if err != nil {
		return nil, fmt.Errorf("failed to open index: %w", err)
	}
	defer index.Close()

	s.index = index

	// Parse and execute query
	q := s.buildQuery(queryStr)
	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Highlight = bleve.NewHighlight()
	searchRequest.Fields = []string{"title", "path", "type", "content"}
	searchRequest.Size = 50

	searchResult, err := index.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Convert results
	results := make([]*SearchResult, 0, len(searchResult.Hits))
	for _, hit := range searchResult.Hits {
		result := &SearchResult{
			ID:    hit.ID,
			Score: hit.Score,
		}

		// Extract fields
		if title, ok := hit.Fields["title"].(string); ok {
			result.Title = title
		}
		if path, ok := hit.Fields["path"].(string); ok {
			result.Path = path
		}
		if fileType, ok := hit.Fields["type"].(string); ok {
			result.Type = fileType
		}

		// Create snippet from highlights or content
		if len(hit.Fragments) > 0 {
			result.Snippet = strings.Join(hit.Fragments["content"], "...")
		} else if content, ok := hit.Fields["content"].(string); ok {
			result.Snippet = s.createSnippet(content, 200)
		}

		results = append(results, result)
	}

	return results, nil
}

// buildQuery builds a Bleve query from the query string
func (s *Searcher) buildQuery(queryStr string) query.Query {
	// Simple query parsing - can be enhanced with KQL later
	queryStr = strings.TrimSpace(queryStr)

	// Check for field-specific queries
	if strings.Contains(queryStr, ":") {
		// Handle field queries like "title:golang" or "type:markdown"
		parts := strings.SplitN(queryStr, ":", 2)
		if len(parts) == 2 {
			field := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			q := bleve.NewMatchQuery(value)
			q.SetField(field)
			return q
		}
	}

	// Check for phrase queries
	if strings.HasPrefix(queryStr, "\"") && strings.HasSuffix(queryStr, "\"") {
		phrase := strings.Trim(queryStr, "\"")
		return bleve.NewMatchPhraseQuery(phrase)
	}

	// Default to match query across all fields
	return bleve.NewMatchQuery(queryStr)
}

// createSnippet creates a snippet from content
func (s *Searcher) createSnippet(content string, maxLength int) string {
	if len(content) <= maxLength {
		return content
	}

	// Try to break at word boundary
	snippet := content[:maxLength]
	if lastSpace := strings.LastIndex(snippet, " "); lastSpace > maxLength/2 {
		snippet = snippet[:lastSpace]
	}

	return snippet + "..."
}

// PrintResults prints search results to console
func PrintResults(results []*SearchResult) {
	if len(results) == 0 {
		fmt.Println("🔍 No results found")
		return
	}

	fmt.Printf("🎯 Found %d results:\n\n", len(results))

	for i, result := range results {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   📁 %s (%s)\n", result.Path, result.Type)
		fmt.Printf("   📊 Score: %.3f\n", result.Score)
		if result.Snippet != "" {
			fmt.Printf("   📝 %s\n", result.Snippet)
		}
		fmt.Println()
	}
}

// SearchByType searches for documents of a specific type
func (s *Searcher) SearchByType(fileType string) ([]*SearchResult, error) {
	query := fmt.Sprintf("type:%s", fileType)
	return s.Search(query)
}

// SearchInContent searches within document content
func (s *Searcher) SearchInContent(content string) ([]*SearchResult, error) {
	query := fmt.Sprintf("content:%s", content)
	return s.Search(query)
}

// SearchByTitle searches in document titles
func (s *Searcher) SearchByTitle(title string) ([]*SearchResult, error) {
	query := fmt.Sprintf("title:%s", title)
	return s.Search(query)
}
