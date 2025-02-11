package realtime

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/DavisAndn/go-aws-crawler/internal/backend"
	"github.com/DavisAndn/go-aws-crawler/internal/config"
)

// Event represents a generic AWS event payload.
type Event struct {
	Detail map[string]interface{} `json:"detail"`
	// Extend with additional fields if needed.
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}
	var event Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	log.Printf("Received event: %+v", event)

	// For now, simply forward the event payload to the backend.
	cfg, err := config.LoadConfig()
	if err == nil {
		_ = backend.SendInitialCrawlResults(cfg.BackendEndpoint, body)
	}

	w.WriteHeader(http.StatusOK)
}

// StartEventServer starts an HTTP server to receive AWS event notifications.
func StartEventServer(addr string) {
	http.HandleFunc("/events", eventHandler)
	log.Printf("Starting event server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
