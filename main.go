package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// The usage of counter metric - amount of valid and invalid authentication
const secret = "NeverGonnaGiveYouApNeverGonnaLetYouDown"

func init() {
	// Setting up random number generator.
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Counter metrics.
	successfulLogin := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "successful_login",
		Help: "Counter representing successful logins to the server",
	})
	invalidLogin := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "invalid_login",
		Help: "Counter representing invalid login attempts to the server",
	})
	// Gauge metrics.
	openConnections := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "open_connections",
		Help: "Number of open connections to the server",
	})
	// Histogram metrics.
	requestSizes := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "request_size_bytes",
		Help:    "Histogram of request sizes for requests to the server",
		Buckets: []float64{100, 200, 300, 400, 500, 1000},
	})

	prometheus.MustRegister(successfulLogin)
	prometheus.MustRegister(invalidLogin)
	prometheus.MustRegister(openConnections)
	prometheus.MustRegister(requestSizes)

	// Exporting metrics registered by the server.
	http.Handle("/metrics", promhttp.Handler())

	/**
	Example of counter metrics.

	Whenever a user hits the "/login" endpoint we simulate a login attempt to the application.
	Corresponding events (success or fail) are reflected in counter metric changes.

	This way (alongside with the usage of PromQL) we can monitor user traffic to the app
	as well as malicious activities happening to our server (huge spike in invalid
	login attempts may be an attack attempt).
	*/
	http.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if value := string(body); value != secret {
			w.WriteHeader(http.StatusForbidden)
			if _, err := w.Write([]byte("Invalid credentials")); err != nil {
				println("Writing response to client:", err)
			}
			invalidLogin.Inc()
			return
		}

		if _, err := w.Write([]byte("You logged in!")); err != nil {
			println("Writing response to client:", err)
		}
		successfulLogin.Inc()
	})

	/**
	Example of gauge metric.

	We count the number of open connections to the "/connect" endpoint. Connections are
	hanging due to random delay between 1-5s.
	*/
	http.HandleFunc("/connect", func(w http.ResponseWriter, req *http.Request) {
		openConnections.Inc()
		defer openConnections.Dec()
		// Simulate some work is being done here.
		delay := rand.Intn(5) + 1
		time.Sleep(time.Duration(delay) * time.Second)
		if _, err := w.Write([]byte(fmt.Sprintf("Your request completed! Your delay was %ds", delay))); err != nil {
			println(err)
		}
	})

	/**
	Example of histogram metric.

	We collect the information about the sizes of requests received by the server.
	*/
	http.HandleFunc("/submit", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		requestSizes.Observe(float64(req.ContentLength))
		if _, err := w.Write([]byte("Thank you for the submission!")); err != nil {
			println(err)
		}
	})

	fmt.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
