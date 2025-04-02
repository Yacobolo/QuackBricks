package client

import (
	"bytes"
	"duckdb-test/cli/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type QueryParam struct {
	Key   string
	Value string
}

// DoAndPrintRequest performs a GET request to the given API path with the provided bearer token
func DoAndPrintRequest(cfg *config.Config, token string, path string, params ...QueryParam) error {
	// Parse base + path as URL
	fullURL, err := url.Parse(fmt.Sprintf("%s%s", cfg.Endpoint, path))
	if err != nil {
		return err
	}

	// Prepare query values
	q := fullURL.Query()
	for _, param := range params {
		q.Add(param.Key, param.Value)
	}
	fullURL.RawQuery = q.Encode() // Encode and attach

	// Build request
	req, err := http.NewRequest("GET", fullURL.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to pretty-print JSON: %w", err)
	}

	fmt.Println(prettyJSON.String())

	return nil
}
