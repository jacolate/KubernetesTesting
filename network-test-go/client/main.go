package main

import (
    "flag"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "time"
)

const (
    defaultServiceURL = "http://h5209.pi.uni-bamberg.de/network-test-go"
    checkInterval    = 100 * time.Millisecond
    timeoutDuration  = 500 * time.Millisecond
    msgCount         = 800
)

func main() {
    // Use the correct time format
    timestamp := time.Now().Format("2006-01-02_15-04-05")
    logFileName := fmt.Sprintf("service_check_%s.log", timestamp)

    logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Error opening log file: %v", err)
    }
    defer logFile.Close()

    multiWriter := io.MultiWriter(os.Stdout, logFile)
    log.SetOutput(multiWriter)

    log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

    serviceURL := flag.String("url", defaultServiceURL, "The URL of the service to check")
    flag.Parse()

    var (
        startUnavailableTime time.Time
        available            = true  // Initialize to true
        failedRequests      int
        completeUnavaiableDuration time.Duration
    )

    client := &http.Client{
        Timeout: timeoutDuration,
    }

    ticker := time.NewTicker(checkInterval)
    defer ticker.Stop()

    counter := 0
    
    log.Printf("Starting service check. Logging to console and %s\n", logFileName)

    for range ticker.C {
        resp, err := client.Get(*serviceURL)

        if err != nil {
            // service is unreachable
            if available {
                startUnavailableTime = time.Now()
                available = false
            }
            log.Printf("Service timeout: %v", err)
            failedRequests++

        } else {
            body, err := ioutil.ReadAll(resp.Body)
            resp.Body.Close()
            if err != nil {
                log.Printf("Error reading response body: %v", err)
            } else {
                // service is reachable again
                if !available {
                    unavailableDuration := time.Since(startUnavailableTime)
                    completeUnavaiableDuration += unavailableDuration
                    log.Printf("Service was unavailable for: %v", unavailableDuration)
                    available = true
                }
                log.Printf("Service is alive: %s", body)
            }
        }
        counter++
        if counter >= msgCount {
            break
        }
    }
    
    log.Printf("Service check completed. Total requests: %d\n", counter)
    log.Printf("Failed requests: %d\n", failedRequests)
    log.Printf("Service was unavailable for a total of: %v", completeUnavaiableDuration)
}
