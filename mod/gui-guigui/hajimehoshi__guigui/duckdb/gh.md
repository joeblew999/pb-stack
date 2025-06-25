# GitHub CLI Tasks

Cross-platform tasks for querying GitHub repositories using the GitHub CLI and gojq.

## Dependencies
- [GitHub CLI](https://cli.github.com/)
- [gojq](https://github.com/itchyny/gojq)

Install both with:
```bash
task gh:deps:install
```

## Common Tasks

List tags:
```bash
task gh:list:tags REPO=owner/repo
```

List latest releases:
```bash
task gh:releases:list REPO=owner/repo LIMIT=5
```

Get repository info:
```bash
task gh:repo:info REPO=owner/repo
```

## Version Management

The tasks support managing versions in a JSON configuration file:

```json
{
  "organization": "owner",
  "version": "1",
  "repositories": [
    {
      "name": "repo1",
      "version": "v1.0.0",
      "description": "Description"
    }
  ]
}
```

Update to latest versions:
```bash
task gh:update:versions CONFIG_FILE=config.json
```

Validate versions:
```bash
task gh:validate:versions CONFIG_FILE=config.json
```

## Links
- [GitHub CLI](https://github.com/cli/cli)
- [gojq](https://github.com/itchyny/gojq)
