module decksh-ebdeck

go 1.23.0

toolchain go1.24.4

require (
	github.com/hajimehoshi/ebiten/v2 v2.9.0-alpha.5.0.20250518103147-cd31850015bb
	github.com/hajimehoshi/guigui v0.0.0-20250606141902-dfd87b3ae4b4
)

require (
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/ebitengine/gomobile v0.0.0-20250329061421-6d0a8e981e4c // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.9.0-alpha.5 // indirect
	github.com/go-text/typesetting v0.3.0 // indirect
	github.com/hajimehoshi/oklab v0.1.0 // indirect
	github.com/jeandeaual/go-locale v0.0.0-20250421151639-a9d6ed1b3d45 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/image v0.28.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	howett.net/plist v1.0.1 // indirect
)

// Note: This module provides a bridge between decksh and ebdeck
// It uses the decksh command externally and implements ebcanvas-style rendering
// The GUI components will be added when network connectivity allows

// Note: decksh and ebcanvas will be used as external commands/libraries
// The integration works by calling the decksh command and implementing
// ebcanvas-style rendering using ebiten directly
