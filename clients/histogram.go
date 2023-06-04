package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func init() {
	// Setting up random number generator.
	rand.Seed(time.Now().UnixNano())
}

func main() {
	client := &http.Client{}
	timer := time.NewTimer(time.Minute)
	// For "attempts" duration time we keep sending the requests to "/connect" endpoint
	// with random delay between 0.5s to 2s.
	// The requests send to "/connect" may be blocking for time between 1s-10s (sleeping on server side).
	for {
		select {
		case <-timer.C:
			return
		default:
			length := int(rand.NormFloat64()*275 + 550)
			if length < 0 {
				length = 0
			}
			// Start the job
			req, err := http.NewRequest("POST", "http://localhost:8080/submit", strings.NewReader(strings.Repeat("a", length)))
			if err != nil {
				fmt.Printf("Error creating request: %v\n", err)
				continue
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error sending request: %v\n", err)
				continue
			}
			_ = resp.Body.Close()
		}
	}
}
