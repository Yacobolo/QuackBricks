package client

import (
	"bytes"
	"duckdb-test/app/internal/config"
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

type RequestParams struct {
	Cfg         *config.Config
	Token       string
	Payload     interface{}
	Path        string
	Method      string
	QueryParams []QueryParam
}

// DoAndPrintRequest performs a GET request to the given API path with the provided bearer token
func DoAndPrintRequest(rp RequestParams) error {
	requestMethod := rp.Method
	if requestMethod == "" {
		requestMethod = http.MethodGet
	}

	// Parse base + path as URL
	fullURL, err := url.Parse(fmt.Sprintf("%s%s", rp.Cfg.Endpoint, rp.Path))
	if err != nil {
		return err
	}

	// Prepare query values
	q := fullURL.Query()
	for _, qp := range rp.QueryParams {
		q.Add(qp.Key, qp.Value)
	}

	fullURL.RawQuery = q.Encode() // Encode and attach

	var reqBody io.Reader
	if rp.Payload != nil {
		jsonData, err := json.Marshal(rp.Payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	// Build request
	req, err := http.NewRequest(requestMethod, fullURL.String(), reqBody)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+rp.Token)
	req.Header.Set("Accept", "application/json")

	// Set Content-Type for methods that have a body
	if rp.Payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body) // Read body for potential error details
		return fmt.Errorf("request failed with status: %s, response: %s", resp.Status, string(bodyBytes))
	}

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, respBodyBytes, "", "  ") // Use 2 spaces for indentation
	if err != nil {
		// If pretty printing fails, just print the raw JSON
		fmt.Println("Failed to pretty-print JSON:", err)
		fmt.Println(string(respBodyBytes))
	} else {
		fmt.Println(prettyJSON.String())
	}

	return nil
}
