package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

const validLoginAttemptsTime = 1 * time.Minute
const invalidLoginAttemptsTime = 10 * time.Second
const invalidLoginAttemptsDelay = 40 * time.Second

func successfulLogin() {
	login(validLoginAttemptsTime, func() time.Duration { return time.Duration(rand.Intn(1500)+500) * time.Millisecond }, "NeverGonnaGiveYouApNeverGonnaLetYouDown")
}

func invalidLogin() {
	time.Sleep(invalidLoginAttemptsDelay)
	login(invalidLoginAttemptsTime, func() time.Duration { return time.Duration(rand.Intn(100)+100) * time.Millisecond }, "hello")
}

func login(attempts time.Duration, reqDelay func() time.Duration, body string) {
	client := &http.Client{}
	timer := time.NewTimer(attempts)
	// For "attempts" duration time we keep sending the requests to "/login" endpoint
	// with random delay between 0.5s to 2s.
	for {
		select {
		case <-timer.C:
			return
		default:
			time.Sleep(reqDelay())
			// Start the job
			req, err := http.NewRequest("POST", "http://localhost:8080/login", strings.NewReader(body))
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

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			successfulLogin()
			wg.Done()
		}()
	}

	wg.Add(1)
	go func() {
		invalidLogin()
		wg.Done()
	}()

	wg.Wait()
}
