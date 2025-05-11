package gui

// Model holds the state for the GUI, like the currently selected view.
type Model struct {
	currentView string
	// Add other global state for the GUI here if needed
}

// Mode returns the current view/mode.
func (m *Model) Mode() string {
	if m.currentView == "" {
		// Default to the first item in the sidebar or a general "home" view
		return "home"
	}
	return m.currentView
}

// SetMode sets the current view/mode.
func (m *Model) SetMode(mode string) {
	m.currentView = mode
}
