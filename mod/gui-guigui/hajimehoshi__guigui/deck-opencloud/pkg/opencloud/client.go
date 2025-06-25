package opencloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client represents an OpenCloud API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
}

// NewClient creates a new OpenCloud client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchRequest represents a search request to OpenCloud
type SearchRequest struct {
	Query     string            `json:"query"`
	Filters   map[string]string `json:"filters,omitempty"`
	Limit     int               `json:"limit,omitempty"`
	Offset    int               `json:"offset,omitempty"`
	SortBy    string            `json:"sort_by,omitempty"`
	SortOrder string            `json:"sort_order,omitempty"`
}

// SearchResponse represents a search response from OpenCloud
type SearchResponse struct {
	Results    []SearchResult `json:"results"`
	Total      int            `json:"total"`
	Took       int            `json:"took"` // milliseconds
	MaxScore   float64        `json:"max_score"`
	Facets     []Facet        `json:"facets,omitempty"`
	Pagination Pagination     `json:"pagination"`
}

// SearchResult represents a single search result
type SearchResult struct {
	ID          string                 `json:"id"`
	Score       float64                `json:"score"`
	Source      map[string]interface{} `json:"source"`
	Highlights  map[string][]string    `json:"highlights,omitempty"`
	Type        string                 `json:"type"`
	Index       string                 `json:"index"`
	Version     int                    `json:"version,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// Facet represents search facets for filtering
type Facet struct {
	Field  string      `json:"field"`
	Values []FacetItem `json:"values"`
}

// FacetItem represents a facet value with count
type FacetItem struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
	Total      int `json:"total"`
}

// IndexRequest represents a document indexing request
type IndexRequest struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Source   map[string]interface{} `json:"source"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// IndexResponse represents an indexing response
type IndexResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Version int    `json:"version"`
	Result  string `json:"result"` // "created" or "updated"
}

// Search performs a search query against OpenCloud
func (c *Client) Search(req *SearchRequest) (*SearchResponse, error) {
	endpoint := fmt.Sprintf("%s/api/search", c.BaseURL)
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search request: %w", err)
	}
	
	httpReq, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	c.setHeaders(httpReq)
	
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}
	
	return &searchResp, nil
}

// Index indexes a document in OpenCloud
func (c *Client) Index(req *IndexRequest) (*IndexResponse, error) {
	endpoint := fmt.Sprintf("%s/api/index", c.BaseURL)
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal index request: %w", err)
	}
	
	httpReq, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	c.setHeaders(httpReq)
	
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("index request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("indexing failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var indexResp IndexResponse
	if err := json.NewDecoder(resp.Body).Decode(&indexResp); err != nil {
		return nil, fmt.Errorf("failed to decode index response: %w", err)
	}
	
	return &indexResp, nil
}

// Delete removes a document from the index
func (c *Client) Delete(id, docType string) error {
	endpoint := fmt.Sprintf("%s/api/index/%s/%s", c.BaseURL, docType, url.PathEscape(id))
	
	httpReq, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}
	
	c.setHeaders(httpReq)
	
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("delete request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	return nil
}

// GetDocument retrieves a document by ID
func (c *Client) GetDocument(id, docType string) (*SearchResult, error) {
	endpoint := fmt.Sprintf("%s/api/document/%s/%s", c.BaseURL, docType, url.PathEscape(id))
	
	httpReq, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get request: %w", err)
	}
	
	c.setHeaders(httpReq)
	
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("get request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("document not found")
	}
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode document: %w", err)
	}
	
	return &result, nil
}

// setHeaders sets common headers for API requests
func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	}
}

// Health checks the health of the OpenCloud service
func (c *Client) Health() error {
	endpoint := fmt.Sprintf("%s/health", c.BaseURL)
	
	resp, err := c.HTTPClient.Get(endpoint)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("service unhealthy: status %d", resp.StatusCode)
	}
	
	return nil
}
