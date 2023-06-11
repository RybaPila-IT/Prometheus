package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

const attempts = 1 * time.Minute
const increaseDelay = 30 * time.Second

func init() {
	// Setting up random number generator.
	rand.Seed(time.Now().UnixNano())
}

func connect() {
	client := &http.Client{}
	timer := time.NewTimer(attempts)
	// For "attempts" duration time we keep sending the requests to "/connect" endpoint
	// with random delay between 0.5s to 2s.
	// The requests send to "/connect" may be blocking for time between 1s-3s (sleeping on server side).
	for {
		select {
		case <-timer.C:
			return
		default:
			time.Sleep(time.Duration(rand.Intn(1500)+500) * time.Millisecond)
			// Start the job
			req, err := http.NewRequest("GET", "http://localhost:8080/connect", strings.NewReader(""))
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

func main() {

	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		go func() {
			connect()
			wg.Done()
		}()
		wg.Add(1)
	}

	time.Sleep(increaseDelay)

	for i := 0; i < 5; i++ {
		go func() {
			connect()
			wg.Done()
		}()
		wg.Add(1)
	}

	wg.Wait()
}
