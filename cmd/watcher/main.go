// cmd/watcher/main.go
package main

import (
	"log"

	"github.com/DavisAndn/go-aws-crawler/internal/realtime"
)

func main() {
	log.Println("Starting watcher event server on port 8282")
	// Start the event server on port 8282.
	realtime.StartEventServer(":8282")
}
