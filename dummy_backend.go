// dummy_backend.go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// handler processes incoming POST requests to /api/aws-resources.
func handler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}

	// Log the received payload
	log.Println("Received data:")
	fmt.Println(string(body))

	// Respond with OK
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Data received successfully")
}

func main() {
	// Set up the route
	http.HandleFunc("/api/aws-resources", handler)

	// Start the HTTP server on port 8181
	log.Println("Dummy backend listening on :8181")
	err := http.ListenAndServe(":8181", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
