package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.Handle("/metrics", VodniSvetKolinMetricsHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
