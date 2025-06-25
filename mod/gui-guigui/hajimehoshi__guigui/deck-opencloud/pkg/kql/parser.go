package kql

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Query represents a parsed KQL query
type Query struct {
	Terms      []Term            `json:"terms"`
	Filters    map[string]string `json:"filters"`
	Sort       []SortField       `json:"sort,omitempty"`
	Limit      int               `json:"limit,omitempty"`
	Offset     int               `json:"offset,omitempty"`
	Aggregates []Aggregate       `json:"aggregates,omitempty"`
}

// Term represents a search term
type Term struct {
	Field    string `json:"field,omitempty"`
	Value    string `json:"value"`
	Operator string `json:"operator"` // "AND", "OR", "NOT"
	Boost    float64 `json:"boost,omitempty"`
	Fuzzy    bool   `json:"fuzzy,omitempty"`
}

// SortField represents a sort specification
type SortField struct {
	Field string `json:"field"`
	Order string `json:"order"` // "asc" or "desc"
}

// Aggregate represents an aggregation request
type Aggregate struct {
	Name  string `json:"name"`
	Field string `json:"field"`
	Type  string `json:"type"` // "terms", "date_histogram", "range", etc.
}

// Parser handles KQL query parsing
type Parser struct {
	defaultOperator string
}

// NewParser creates a new KQL parser
func NewParser() *Parser {
	return &Parser{
		defaultOperator: "AND",
	}
}

// Parse parses a KQL query string into a structured Query
func (p *Parser) Parse(queryStr string) (*Query, error) {
	query := &Query{
		Terms:   []Term{},
		Filters: make(map[string]string),
		Limit:   50, // default limit
	}

	if strings.TrimSpace(queryStr) == "" {
		return query, nil
	}

	// Split query into tokens
	tokens := p.tokenize(queryStr)
	
	for i, token := range tokens {
		if err := p.parseToken(token, query, i, tokens); err != nil {
			return nil, fmt.Errorf("failed to parse token '%s': %w", token, err)
		}
	}

	return query, nil
}

// tokenize splits the query string into tokens
func (p *Parser) tokenize(queryStr string) []string {
	// Simple tokenization - can be enhanced
	// Handle quoted strings, parentheses, operators
	
	// Regex to split on spaces but preserve quoted strings
	re := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)'`)
	matches := re.FindAllString(queryStr, -1)
	
	var tokens []string
	for _, match := range matches {
		// Remove quotes from quoted strings
		if (strings.HasPrefix(match, `"`) && strings.HasSuffix(match, `"`)) ||
		   (strings.HasPrefix(match, `'`) && strings.HasSuffix(match, `'`)) {
			match = match[1 : len(match)-1]
		}
		tokens = append(tokens, match)
	}
	
	return tokens
}

// parseToken parses a single token
func (p *Parser) parseToken(token string, query *Query, index int, allTokens []string) error {
	// Handle special keywords
	switch strings.ToUpper(token) {
	case "AND", "OR", "NOT":
		// Operator tokens are handled when processing terms
		return nil
	case "LIMIT":
		return p.parseLimit(index, allTokens, query)
	case "SORT", "ORDER":
		return p.parseSort(index, allTokens, query)
	}

	// Handle field:value syntax
	if strings.Contains(token, ":") {
		return p.parseFieldValue(token, query)
	}

	// Handle range queries [value TO value]
	if strings.HasPrefix(token, "[") && strings.HasSuffix(token, "]") {
		return p.parseRange(token, query)
	}

	// Regular search term
	term := Term{
		Value:    token,
		Operator: p.defaultOperator,
	}
	
	// Check for fuzzy search (~)
	if strings.HasSuffix(token, "~") {
		term.Value = strings.TrimSuffix(token, "~")
		term.Fuzzy = true
	}
	
	// Check for boost (^number)
	if strings.Contains(token, "^") {
		parts := strings.Split(token, "^")
		if len(parts) == 2 {
			if boost, err := strconv.ParseFloat(parts[1], 64); err == nil {
				term.Value = parts[0]
				term.Boost = boost
			}
		}
	}

	query.Terms = append(query.Terms, term)
	return nil
}

// parseFieldValue parses field:value syntax
func (p *Parser) parseFieldValue(token string, query *Query) error {
	parts := strings.SplitN(token, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid field:value syntax: %s", token)
	}

	field := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	// Handle special fields
	switch field {
	case "type", "content_type", "format":
		query.Filters["type"] = value
	case "size":
		query.Filters["size"] = value
	case "created", "modified", "date":
		query.Filters["date"] = value
	case "author", "creator":
		query.Filters["author"] = value
	case "tag", "tags":
		query.Filters["tags"] = value
	default:
		// Regular field search
		term := Term{
			Field:    field,
			Value:    value,
			Operator: p.defaultOperator,
		}
		query.Terms = append(query.Terms, term)
	}

	return nil
}

// parseRange parses range queries like [1 TO 100] or [2023-01-01 TO 2023-12-31]
func (p *Parser) parseRange(token string, query *Query) error {
	// Remove brackets
	rangeStr := token[1 : len(token)-1]
	
	// Split on TO
	parts := strings.Split(strings.ToUpper(rangeStr), " TO ")
	if len(parts) != 2 {
		return fmt.Errorf("invalid range syntax: %s", token)
	}

	from := strings.TrimSpace(parts[0])
	to := strings.TrimSpace(parts[1])

	// For now, add as a filter - can be enhanced
	query.Filters["range"] = fmt.Sprintf("%s TO %s", from, to)
	
	return nil
}

// parseLimit parses LIMIT clause
func (p *Parser) parseLimit(index int, tokens []string, query *Query) error {
	if index+1 >= len(tokens) {
		return fmt.Errorf("LIMIT requires a number")
	}

	limitStr := tokens[index+1]
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return fmt.Errorf("invalid LIMIT value: %s", limitStr)
	}

	query.Limit = limit
	return nil
}

// parseSort parses SORT/ORDER BY clause
func (p *Parser) parseSort(index int, tokens []string, query *Query) error {
	if index+1 >= len(tokens) {
		return fmt.Errorf("SORT requires a field name")
	}

	field := tokens[index+1]
	order := "asc" // default

	// Check for explicit order
	if index+2 < len(tokens) {
		orderToken := strings.ToLower(tokens[index+2])
		if orderToken == "asc" || orderToken == "desc" {
			order = orderToken
		}
	}

	sortField := SortField{
		Field: field,
		Order: order,
	}

	query.Sort = append(query.Sort, sortField)
	return nil
}

// ToBleve converts KQL query to Bleve query format
func (q *Query) ToBleve() map[string]interface{} {
	bleveQuery := make(map[string]interface{})
	
	if len(q.Terms) > 0 {
		if len(q.Terms) == 1 {
			term := q.Terms[0]
			if term.Field != "" {
				bleveQuery["field"] = term.Field
			}
			bleveQuery["match"] = term.Value
		} else {
			// Multiple terms - create boolean query
			should := make([]map[string]interface{}, 0, len(q.Terms))
			for _, term := range q.Terms {
				termQuery := map[string]interface{}{
					"match": term.Value,
				}
				if term.Field != "" {
					termQuery["field"] = term.Field
				}
				should = append(should, termQuery)
			}
			bleveQuery["bool"] = map[string]interface{}{
				"should": should,
			}
		}
	}

	return bleveQuery
}

// ToOpenCloud converts KQL query to OpenCloud search request format
func (q *Query) ToOpenCloud() map[string]interface{} {
	request := make(map[string]interface{})
	
	// Build query string
	var queryParts []string
	for _, term := range q.Terms {
		if term.Field != "" {
			queryParts = append(queryParts, fmt.Sprintf("%s:%s", term.Field, term.Value))
		} else {
			queryParts = append(queryParts, term.Value)
		}
	}
	
	if len(queryParts) > 0 {
		request["query"] = strings.Join(queryParts, " ")
	}
	
	// Add filters
	if len(q.Filters) > 0 {
		request["filters"] = q.Filters
	}
	
	// Add sorting
	if len(q.Sort) > 0 {
		sort := make([]map[string]string, 0, len(q.Sort))
		for _, s := range q.Sort {
			sort = append(sort, map[string]string{
				"field": s.Field,
				"order": s.Order,
			})
		}
		request["sort"] = sort
	}
	
	// Add pagination
	if q.Limit > 0 {
		request["limit"] = q.Limit
	}
	if q.Offset > 0 {
		request["offset"] = q.Offset
	}
	
	return request
}

// String returns a string representation of the query
func (q *Query) String() string {
	var parts []string
	
	for _, term := range q.Terms {
		if term.Field != "" {
			parts = append(parts, fmt.Sprintf("%s:%s", term.Field, term.Value))
		} else {
			parts = append(parts, term.Value)
		}
	}
	
	result := strings.Join(parts, " ")
	
	if len(q.Sort) > 0 {
		result += " SORT " + q.Sort[0].Field + " " + q.Sort[0].Order
	}
	
	if q.Limit > 0 {
		result += fmt.Sprintf(" LIMIT %d", q.Limit)
	}
	
	return result
}
