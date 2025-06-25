package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

type CommandRequest struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type CommandResponse struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) IsAvailable() bool {
	resp, err := c.client.Get(c.baseURL + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK
}

func (c *Client) ExecuteCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command specified")
	}

	command := args[0]
	commandArgs := args[1:]

	req := CommandRequest{
		Command: command,
		Args:    commandArgs,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.client.Post(
		c.baseURL+"/api/v1/tasks",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var response CommandResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("command failed: %s", response.Error)
	}

	fmt.Print(response.Output)
	return nil
}

func (c *Client) Which(binary string) (string, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/api/v1/tools/which/%s", c.baseURL, binary))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Found bool   `json:"found"`
		Path  string `json:"path"`
		Error string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Found {
		return "", fmt.Errorf("binary not found: %s", response.Error)
	}

	return response.Path, nil
}

func (c *Client) Download(url, output string) error {
	req := struct {
		URL    string `json:"url"`
		Output string `json:"output"`
	}{
		URL:    url,
		Output: output,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.client.Post(
		c.baseURL+"/api/v1/tools/got",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("download failed: %s", response.Error)
	}

	return nil
}
