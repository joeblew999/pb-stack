# Font Manifest

Generated: 2025-06-16 00:13:40
Total fonts: 5

## Available Fonts

| Font File | Display Name | Family | Weight | Style |
|-----------|--------------|--------|--------|
| `FiraSans-Regular.ttf` | Fira Sans Regular | Fira Sans | 400 | normal |
| `FiraSans-Bold.ttf` | Fira Sans Bold | Fira Sans | 700 | normal |
| `FiraMono-Regular.ttf` | Fira Mono Regular | Fira Mono | 400 | normal |
| `Lato-Regular.ttf` | Lato Regular | Lato | 400 | normal |
| `LibreBaskerville-Regular.ttf` | Libre Baskerville Regular | Libre Baskerville | 400 | normal |

## Usage in decksh

Set the font directory:
```bash
export DECKFONTS="$(pwd)/fonts"
# or use -fontdir flag
pngdeck -fontdir ./fonts input.xml output.png
```

## Font Notes

- Some fonts are downloaded as WOFF2 format (web fonts)
- decksh works best with TTF/OTF fonts
- Consider converting WOFF2 to TTF if needed
- Fira Sans and Fira Mono are reliably available as TTF
