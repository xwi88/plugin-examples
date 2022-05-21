package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// http exporter
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":12022", nil))
}
