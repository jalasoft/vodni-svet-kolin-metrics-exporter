package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var port uint64

func init() {
	flag.Uint64Var(&port, "port", 8080, "Server port")
}

func main() {

	log.SetOutput(os.Stdout)

	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle("/metrics", VodniSvetKolinMetricsHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

	log.Printf("Server started on port %d", port)
}
