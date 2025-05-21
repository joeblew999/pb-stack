package assettypes

// --- GitHub Release Asset Finder Types ---

// AssetFilters defines the filtering criteria for assets.
type AssetFilters struct {
	VersionContains string `yaml:"versionContains,omitempty"`
	OsContains      string `yaml:"osContains,omitempty"`
	ArchContains    string `yaml:"archContains,omitempty"`
	ExtContains     string `yaml:"extContains,omitempty"`
	NameContains    string `yaml:"nameContains,omitempty"`
}

// AssetSearchEntry defines a single search configuration.
type AssetSearchEntry struct {
	Name         string       `yaml:"name"`                  // For identifying the search in output
	Description  string       `yaml:"description,omitempty"` // Optional description
	Owner        string       `yaml:"owner"`
	Repo         string       `yaml:"repo"`
	Tag          string       `yaml:"tag"` // Can also be "latest"
	Filters      AssetFilters `yaml:"filters"`
	SelectField  string       `yaml:"selectField,omitempty"`  // Field to print if not downloading
	FirstOnly    bool         `yaml:"firstOnly,omitempty"`    // Default to true?
	Download     bool         `yaml:"download,omitempty"`     // Whether to download the found asset
	DownloadPath string       `yaml:"downloadPath,omitempty"` // Optional path to download to. If empty, uses current dir + asset name.
}

// AssetSearchConfig is the top-level structure for the YAML config file.
type AssetSearchConfig struct {
	Searches []AssetSearchEntry `yaml:"searches"`
}

// GitHubRelease represents the structure of a GitHub release API response.
type GitHubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []GitHubAsset `json:"assets"`
	Message string        `json:"message"` // For API error messages like "Not Found"
}

// GitHubAsset represents a single asset within a GitHub release.
type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int    `json:"size"`
	ContentType        string `json:"content_type"`
}

// --- End GitHub Release Asset Finder Types ---
