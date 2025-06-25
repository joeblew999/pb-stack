package indexer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Document represents an indexed document
type Document struct {
	ID       string                 `json:"id"`
	Title    string                 `json:"title"`
	Content  string                 `json:"content"`
	Path     string                 `json:"path"`
	Type     string                 `json:"type"`
	Size     int64                  `json:"size"`
	ModTime  time.Time              `json:"mod_time"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Tags     []string               `json:"tags,omitempty"`
	Language string                 `json:"language,omitempty"`
}

// Indexer handles document indexing
type Indexer struct {
	index    bleve.Index
	indexDir string
}

// New creates a new indexer
func New(indexDir string) *Indexer {
	return &Indexer{
		indexDir: indexDir,
	}
}

// IndexDirectory indexes all supported files in a directory
func (idx *Indexer) IndexDirectory(dataDir string) error {
	// Create or open index
	if err := idx.openOrCreateIndex(); err != nil {
		return fmt.Errorf("failed to open index: %w", err)
	}
	defer idx.index.Close()

	fmt.Printf("📁 Scanning directory: %s\n", dataDir)

	return filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file should be indexed
		if idx.shouldIndex(path) {
			fmt.Printf("📄 Indexing: %s\n", path)
			if err := idx.indexFile(path, info); err != nil {
				fmt.Printf("⚠️  Failed to index %s: %v\n", path, err)
			}
		}

		return nil
	})
}

// shouldIndex determines if a file should be indexed
func (idx *Indexer) shouldIndex(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	supportedExts := []string{
		".md", ".markdown", // Markdown files
		".txt",       // Text files
		".go",        // Go source files
		".js", ".ts", // JavaScript/TypeScript
		".py",           // Python
		".yaml", ".yml", // YAML files
		".json",         // JSON files
		".html", ".htm", // HTML files
	}

	for _, supported := range supportedExts {
		if ext == supported {
			return true
		}
	}
	return false
}

// indexFile indexes a single file
func (idx *Indexer) indexFile(path string, info os.FileInfo) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	doc := &Document{
		ID:      path,
		Path:    path,
		Type:    idx.getFileType(path),
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}

	// Process based on file type
	switch doc.Type {
	case "markdown":
		if err := idx.processMarkdown(doc, content); err != nil {
			return err
		}
	default:
		doc.Content = string(content)
		doc.Title = filepath.Base(path)
	}

	// Index the document
	return idx.index.Index(doc.ID, doc)
}

// processMarkdown processes markdown files with metadata extraction
func (idx *Indexer) processMarkdown(doc *Document, content []byte) error {
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	context := parser.NewContext()
	mdDoc := md.Parser().Parse(text.NewReader(content), parser.WithContext(context))

	// Extract metadata
	metaData := meta.Get(context)
	if metaData != nil {
		doc.Metadata = metaData

		// Extract title from metadata or use filename
		if title, ok := metaData["title"].(string); ok {
			doc.Title = title
		} else {
			doc.Title = strings.TrimSuffix(filepath.Base(doc.Path), filepath.Ext(doc.Path))
		}

		// Extract tags
		if tags, ok := metaData["tags"].([]interface{}); ok {
			for _, tag := range tags {
				if tagStr, ok := tag.(string); ok {
					doc.Tags = append(doc.Tags, tagStr)
				}
			}
		}

		// Extract language
		if lang, ok := metaData["language"].(string); ok {
			doc.Language = lang
		}
	} else {
		doc.Title = strings.TrimSuffix(filepath.Base(doc.Path), filepath.Ext(doc.Path))
	}

	// Convert markdown to plain text for indexing
	var buf strings.Builder
	if err := md.Renderer().Render(&buf, content, mdDoc); err != nil {
		return err
	}

	// Store both original markdown and rendered content
	doc.Content = string(content)

	return nil
}

// getFileType determines the file type based on extension
func (idx *Indexer) getFileType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".md", ".markdown":
		return "markdown"
	case ".go":
		return "go"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".py":
		return "python"
	case ".yaml", ".yml":
		return "yaml"
	case ".json":
		return "json"
	case ".html", ".htm":
		return "html"
	default:
		return "text"
	}
}

// openOrCreateIndex opens existing index or creates new one
func (idx *Indexer) openOrCreateIndex() error {
	index, err := bleve.Open(idx.indexDir)
	if err == bleve.ErrorIndexPathDoesNotExist {
		fmt.Println("📚 Creating new index...")
		mapping := idx.createIndexMapping()
		index, err = bleve.New(idx.indexDir, mapping)
	}

	if err != nil {
		return err
	}

	idx.index = index
	return nil
}

// createIndexMapping creates the index mapping for documents
func (idx *Indexer) createIndexMapping() mapping.IndexMapping {
	// Create document mapping
	docMapping := bleve.NewDocumentMapping()

	// Text fields
	textFieldMapping := bleve.NewTextFieldMapping()
	docMapping.AddFieldMappingsAt("title", textFieldMapping)
	docMapping.AddFieldMappingsAt("content", textFieldMapping)
	docMapping.AddFieldMappingsAt("path", textFieldMapping)

	// Keyword fields
	keywordFieldMapping := bleve.NewKeywordFieldMapping()
	docMapping.AddFieldMappingsAt("type", keywordFieldMapping)
	docMapping.AddFieldMappingsAt("language", keywordFieldMapping)

	// Date field
	dateFieldMapping := bleve.NewDateTimeFieldMapping()
	docMapping.AddFieldMappingsAt("mod_time", dateFieldMapping)

	// Numeric field
	numericFieldMapping := bleve.NewNumericFieldMapping()
	docMapping.AddFieldMappingsAt("size", numericFieldMapping)

	// Index mapping
	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("_default", docMapping)

	return indexMapping
}
