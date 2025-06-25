package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"opencloud/pkg/indexer" // Assuming this path is correct
	"opencloud/pkg/search"  // Assuming this path is correct
)

var (
	mode     = flag.String("mode", "", "Mode: index, search (required)")
	indexDir = flag.String("index", "./index", "Index directory")
	dataDir  = flag.String("data", "./data", "Data directory to index (for index mode)")
	query    = flag.String("query", "", "Search query (for search mode)")
)

func main() {
	flag.Parse()

	if *mode == "" {
		fmt.Println("❌ Mode (-mode) is required: 'index' or 'search'")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("🌩️  OpenCloud - CLI Mode")
	fmt.Println("========================")

	switch *mode {
	case "index":
		fmt.Printf("📚 Indexing directory: %s into %s\n", *dataDir, *indexDir)
		idx := indexer.New(*indexDir) // Renamed 'indexer' to 'idx' to avoid conflict with package name
		if err := idx.IndexDirectory(*dataDir); err != nil {
			log.Fatalf("Indexing failed: %v", err)
		}
		fmt.Println("✅ Indexing completed")

	case "search":
		if *query == "" {
			fmt.Println("❌ Query (-query) required for search mode")
			flag.Usage()
			os.Exit(1)
		}
		fmt.Printf("🔍 Searching for: \"%s\" in index: %s\n", *query, *indexDir)
		searcher := search.New(*indexDir)
		results, err := searcher.Search(*query)
		if err != nil {
			log.Fatalf("Search failed: %v", err)
		}
		search.PrintResults(results)

	default:
		fmt.Printf("❌ Unknown mode: %s. Available modes: index, search.\n", *mode)
		flag.Usage()
		os.Exit(1)
	}
}
