# Git Project Management

Cross-platform tasks for managing multiple Git repositories with version control.

## Configuration

Uses a JSON configuration file (default: `repos.json`):
```json
{
  "organization": "org-name",
  "version": "1",
  "repositories": [
    {
      "name": "repo-name",
      "version": "v1.0.0",
      "description": "Repository description"
    }
  ]
}
```

## Common Tasks

Clone all repositories:
```bash
task git:clone:all [CONFIG_FILE=repos.json] [BUILD_DIR=build]
```

Update all repositories to specified versions:
```bash
task git:update:all [CONFIG_FILE=repos.json] [BUILD_DIR=build]
```

Clean all repositories:
```bash
task git:clean:all [CONFIG_FILE=repos.json] [BUILD_DIR=build]
```

## Version Management

Update config with latest versions and sync repositories:
```bash
task git:sync:versions [CONFIG_FILE=repos.json]
```

Validate versions and repository states:
```bash
task git:check:versions [CONFIG_FILE=repos.json]
```

## Dependencies
- Git
- [GitHub CLI](https://cli.github.com/)
- [gojq](https://github.com/itchyny/gojq)

Install dependencies:
```bash
task git:deps:install
```
