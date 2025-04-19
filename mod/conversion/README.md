# conversion

Pipelines for conversions so that:

AI can get clean markdown

Users can get clean markdown



## HTML

https://github.com/mackee/go-readability

Can call remote or local html.

remote sites often block.


```sh
# Install the CLI tool
go install github.com/mackee/go-readability/cmd/readability@latest

# Extract content from a URL
readability https://example.com/article

# Save the extracted content to a file
readability https://example.com/article > article.html

# Output as markdown
readability --format markdown https://example.com/article > article.md

# Output metadata as JSON
readability --metadata https://example.com/article
```

## DATA



