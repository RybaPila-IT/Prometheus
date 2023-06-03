package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// The usage of counter metric - amount of valid and invalid authentication
const secret = "NeverGonnaGiveYouApNeverGonnaLetYouDown"

func main() {

	successfulLogin := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "successful_login",
		Help: "Counter representing successful logins to the application",
	})
	invalidLogin := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "invalid_login",
		Help: "Counter representing invalid login attempts to the application",
	})

	prometheus.MustRegister(successfulLogin)
	prometheus.MustRegister(invalidLogin)

	http.Handle("/metrics", promhttp.Handler())

	/**

	Example of counter metrics.

	Whenever a user hits the "/login" endpoint we simulate a login attempt to the application.
	Corresponding events (success or fail) are reflected in counter metric changes.

	This way (alongside with the usage of PromQL) we can monitor user traffic to the app
	as well as malicious activities happening to our server (huge spike in invalid
	login attempts may be an attack attempt).

	*/
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
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

	fmt.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
