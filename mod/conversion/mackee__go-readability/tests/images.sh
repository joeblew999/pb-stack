# Install the CLI tool
go install github.com/mackee/go-readability/cmd/readability@latest

# Extract content from a URL
readability https://unsplash.com/s/photos/sexy-woman

mkdir -p out

# Save the extracted content to a file
readability https://unsplash.com/s/photos/sexy-woman > out/images.html

# Output as markdown
readability --format markdown https://unsplash.com/s/photos/sexy-woman > out/images.md

# Output metadata as JSON
readability --metadata https://unsplash.com/s/photos/sexy-woman > out/images.json