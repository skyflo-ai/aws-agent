package backend

import (
	"bytes"
	"fmt"
	"net/http"
)

// SendInitialCrawlResults posts the JSON payload to the backend endpoint.
func SendInitialCrawlResults(backendURL string, payload []byte) error {
	req, err := http.NewRequest("POST", backendURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("backend returned status %d", resp.StatusCode)
	}

	return nil
}
