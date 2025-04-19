# Install the CLI tool
go install github.com/mackee/go-readability/cmd/readability@latest

# Extract content from a URL
readability https://data-star.dev/examples/click_to_edit 

mkdir -p out

# Save the extracted content to a file
readability https://data-star.dev/examples/click_to_edit > out/datastar.html

# Output as markdown
readability --format markdown https://data-star.dev/examples/click_to_edit > out/datastar.md

# Output metadata as JSON
readability --metadata https://data-star.dev/examples/click_to_edit > out/datastar.json