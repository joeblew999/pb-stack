# conversion

These are basic workers used in pipelines to convert from one format to another.

This helps with Single Sourcing so that we output to many other system like Google SEO, etc. 

It has a myriad of uses at different levels of the Software Flow though.

## Markdown

Pipelines for conversions so that:

AI can get clean markdown

Users can get clean markdown


## HTML

Convert complex web pages to simple ones, to aid with AI and SEO.

https://github.com/mackee/go-readability

Can call remote or local html.

Remote sites often block. 


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



