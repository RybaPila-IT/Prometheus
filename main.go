package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "my_counter",
		Help: "This is my counter metric",
	})
	prometheus.MustRegister(counter)

	go func() {
		for {
			counter.Inc()
			time.Sleep(2 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	// Start the web server
	fmt.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
