package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"main/pkg/assettypes" // Import the new package

	"gopkg.in/yaml.v3"
)

// --- GitHub Release Asset Finder ---

// FindAssetsFromConfig reads a YAML config file and processes each search entry.
func FindAssetsFromConfig(configFilePath string, apiToken string) {
	log.Printf("Loading asset search configurations from: %s", configFilePath)
	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Error reading asset config file %s: %v", configFilePath, err)
	}

	var config assettypes.AssetSearchConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error unmarshaling asset config YAML from %s: %v", configFilePath, err)
	}

	if len(config.Searches) == 0 {
		log.Println("No searches defined in the configuration file.")
		return
	}

	log.Printf("Found %d search configuration(s).", len(config.Searches))
	for i, search := range config.Searches {
		log.Printf("--- Processing search #%d: '%s' ---", i+1, search.Name)
		if search.Description != "" {
			log.Printf("    Description: %s", search.Description)
		}

		findSingleGitHubReleaseAsset(
			search.Owner, search.Repo, search.Tag, apiToken,
			search.Filters.VersionContains, search.Filters.OsContains, search.Filters.ArchContains,
			search.Filters.ExtContains, search.Filters.NameContains, search.SelectField,
			search.FirstOnly, search.Download, search.DownloadPath,
		)
	}
	log.Println("--- All asset searches from config complete ---")
}

// findSingleGitHubReleaseAsset is the core logic for finding an asset based on parameters.
func findSingleGitHubReleaseAsset(owner, repo, tag, token, versionContains, osContains, archContains, extContains, nameContains, selectField string, firstOnly bool, download bool, downloadPath string) {
	if owner == "" || repo == "" || tag == "" {
		log.Println("Error: Owner, repo, and tag are required for finding an asset.")
		return
	}

	apiURL := buildGitHubAPIURL(owner, repo, tag)
	log.Printf("Fetching release data from: %s", apiURL)
	releaseData, err := fetchReleaseDataFromAPI(apiURL, token)
	if err != nil {
		log.Fatalf("Error fetching release data: %v", err)
	}

	if releaseData.Message != "" && len(releaseData.Assets) == 0 {
		log.Fatalf("GitHub API returned a message (likely an error): %s for tag %s", releaseData.Message, tag)
	}
	if len(releaseData.Assets) == 0 {
		log.Printf("No assets found for release tag '%s'.", tag)
		return
	}

	var matchingAssets []assettypes.GitHubAsset
	for _, asset := range releaseData.Assets {
		if filterGitHubAsset(asset, versionContains, osContains, archContains, extContains, nameContains) {
			matchingAssets = append(matchingAssets, asset)
		}
	}

	if len(matchingAssets) == 0 {
		log.Println("No assets matched the filter criteria.")
		return
	}

	assetsToProcess := matchingAssets
	if firstOnly && len(matchingAssets) > 0 {
		assetsToProcess = []assettypes.GitHubAsset{matchingAssets[0]}
	}

	if download {
		log.Printf("Download requested. Processing %d matching asset(s).", len(assetsToProcess))
		for _, asset := range assetsToProcess {
			targetPath := downloadPath
			if targetPath == "" {
				targetPath = asset.Name
			}
			log.Printf("Downloading '%s' to '%s'...", asset.BrowserDownloadURL, targetPath)
			if err := downloadFile(asset.BrowserDownloadURL, targetPath); err != nil {
				log.Printf("Error downloading asset '%s': %v", asset.Name, err)
			}
		}
	} else {
		for _, asset := range assetsToProcess {
			printSelectedAssetField(asset, selectField)
		}
	}
}

func buildGitHubAPIURL(owner, repo, tag string) string {
	if strings.ToLower(tag) == "latest" {
		return fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	}
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, repo, tag)
}

func fetchReleaseDataFromAPI(apiURL, token string) (assettypes.GitHubRelease, error) {
	var release assettypes.GitHubRelease
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return release, fmt.Errorf("creating request for API %s: %w", apiURL, err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return release, fmt.Errorf("fetching release from API %s: %w", apiURL, err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return release, fmt.Errorf("reading response body from API %s: %w", apiURL, err)
	}

	if err := json.Unmarshal(bodyBytes, &release); err != nil {
		return release, fmt.Errorf("decoding JSON response from API %s (body: %s): %w", apiURL, string(bodyBytes), err)
	}

	if resp.StatusCode != http.StatusOK && release.Message == "" {
		return release, fmt.Errorf("API request to %s failed with status %s: %s", apiURL, resp.Status, string(bodyBytes))
	}

	return release, nil
}

func downloadFile(url string, targetPath string) error {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to perform GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %s", resp.Status)
	}

	targetDir := filepath.Dir(targetPath)
	if targetDir != "." {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create target directory '%s': %w", targetDir, err)
		}
	}

	outFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", targetPath, err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}

func filterGitHubAsset(asset assettypes.GitHubAsset, version, osFilter, arch, ext, name string) bool {
	assetNameLower := strings.ToLower(asset.Name)
	if version != "" && !strings.Contains(asset.Name, version) {
		return false
	}
	if osFilter != "" && !strings.Contains(assetNameLower, strings.ToLower(osFilter)) {
		return false
	}
	if arch != "" && !strings.Contains(assetNameLower, strings.ToLower(arch)) {
		return false
	}
	if ext != "" && !strings.HasSuffix(assetNameLower, strings.ToLower(ext)) {
		return false
	}
	if name != "" && !strings.Contains(assetNameLower, strings.ToLower(name)) {
		return false
	}
	return true
}

func printSelectedAssetField(asset assettypes.GitHubAsset, field string) {
	switch strings.ToLower(field) {
	case "name":
		fmt.Println(asset.Name)
	case "browser_download_url":
		fmt.Println(asset.BrowserDownloadURL)
	case "size":
		fmt.Println(asset.Size)
	case "content_type":
		fmt.Println(asset.ContentType)
	default:
		log.Printf("Unknown field to select: '%s'. Defaulting to browser_download_url.", field)
		fmt.Println(asset.BrowserDownloadURL)
	}
}
