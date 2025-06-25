# Card Game in Go
## Project Structure

```
├── cmd/            # Entry points (cardgen, game, fontdownloader)
├── .assets/        # All generated assets (dot-prefixed = not committed)
│   ├── xml/        # Generated XML from decksh
│   ├── svg/        # Generated SVG files
│   ├── png/        # Generated PNG files for game
│   └── fonts/      # Downloaded Google Fonts
├── fonts.json      # Font configuration file
└── .bin/           # Local Go tools (svgdeck, pngdeck)
```ard game using Ebiten for graphics and decksh/deck for high-quality card generation.

## Quick Start

```bash
# Install tools and fonts
task dev:install-deck-tools
task dev:install-fonts

# Generate cards
task cards:decksh-png

# Run the game
task game:run
```

## Project Structure

```
├── cmd/            # Entry points (cardgen, game, fontdownloader)
├── .assets/        # All generated assets (not committed)
│   ├── xml/        # Generated XML from decksh
│   ├── svg/        # Generated SVG files
│   ├── png/        # Generated PNG files for game
│   └── fonts/      # Downloaded Google Fonts
└── .bin/           # Local Go tools (svgdeck, pngdeck)
```

## Card Generation

Uses [decksh](https://github.com/ajstarks/decksh) to generate professional card designs:

- **decksh** → XML markup
- **svgdeck** → SVG graphics  
- **pngdeck** → PNG images

All tools are installed locally, no external dependencies required.

**Font Management**: Fonts are configured via `fonts.json` and downloaded automatically.

## Fonts downloading

https://fonts.google.com



## Available Tasks

```bash
task --list
```

Key tasks:
- `cards:decksh-png` - Generate card PNGs
- `game:run` - Start the game
- `clean:all` - Remove all generated files

## Font Configuration

Fonts are managed via `fonts.json`. To add a new font, edit the file:

```json
{
  "fonts": [
    {
      "name": "FontName-Regular.ttf",
      "displayName": "Font Display Name", 
      "url": "https://github.com/google/fonts/raw/main/path/to/font.ttf",
      "family": "Font Family",
      "weight": "400",
      "style": "normal"
    }
  ]
}
```

Run `task fonts:gen` to download any new fonts.

## Font Management

### JSON-Driven Configuration

Fonts are managed via `fonts.json` configuration:

```json
{
  "fonts": [
    {
      "name": "FiraSans-Regular.ttf",
      "displayName": "Fira Sans Regular", 
      "family": "Fira Sans",
      "weight": "400",
      "style": "normal",
      "url": "https://github.com/google/fonts/raw/main/ofl/firasans/FiraSans-Regular.ttf"
    }
  ]
}
```

### Enhanced Font Discovery (Optional)

The fontdownloader now supports Google Fonts API integration:

```bash
# Search for fonts
go run ./cmd/fontdownloader --search "Roboto"

# List fonts by category
go run ./cmd/fontdownloader --list --category serif

# Add a font to your project
go run ./cmd/fontdownloader --add "Open Sans"

# Update existing fonts with latest URLs
go run ./cmd/fontdownloader --update
```

**Note**: Google Fonts API requires an API key for search/discovery features. The basic download functionality works without an API key.

To enable full functionality:
```bash
export GOOGLE_FONTS_API_KEY="your_api_key_here"
go run ./cmd/fontdownloader --api-key "$GOOGLE_FONTS_API_KEY" --search "Roboto"
```

Get an API key at: https://developers.google.com/fonts/docs/developer_api




