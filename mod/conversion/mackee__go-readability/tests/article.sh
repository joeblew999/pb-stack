# Install the CLI tool
go install github.com/mackee/go-readability/cmd/readability@latest

readability -h 

# Extract content from a URL
readability https://example.com/article

# Save the extracted content to a file
readability https://example.com/article > out/article.html

# Output as markdown
readability --format markdown https://example.com/article > out/article.md

# Output metadata as JSON
readability --metadata https://example.com/article out/article.json

