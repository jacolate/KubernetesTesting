package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	defaultServiceURL = "http://h5209.pi.uni-bamberg.de/network-test-go" // Replace with your default service URL
	checkInterval     = 1 * time.Millisecond
	timeoutDuration   = 1 * time.Second
)

func main() {
	// Define the command-line flag for the service URL
	serviceURL := flag.String("url", defaultServiceURL, "The URL of the service to check")
	flag.Parse()

	var (
		startUnavailableTime time.Time
		isUnavailable        bool
	)

	client := &http.Client{
		Timeout: timeoutDuration,
	}

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		resp, err := client.Get(*serviceURL)
		if err != nil {
			if !isUnavailable {
				startUnavailableTime = time.Now()
				isUnavailable = true
			}
			log.Printf("Service timeout: %v", err)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Printf("Error reading response body: %v", err)
			} else {
				if isUnavailable {
					unavailableDuration := time.Since(startUnavailableTime)
					log.Printf("Service was unavailable for: %v", unavailableDuration)
					isUnavailable = false
				}
				log.Printf("Service is alive: %s", body)
			}
		}
	}
}
