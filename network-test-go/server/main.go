package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
    "os"
)

func main() {
    nodeName := os.Getenv("NODE_NAME")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now().Format(time.RFC1123)
		w.WriteHeader(http.StatusOK)
        log.Printf("Service accessed from %s\n", r.RemoteAddr)
		fmt.Fprintf(w, "Service reachable, The current time is %s. Running on node: %s\n", currentTime, nodeName)
	})

	port := "8080"
	log.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
