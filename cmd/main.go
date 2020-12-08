package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	port      = flag.String("port", ":8080", "The port")
	staticDir = flag.String("staticDir", "./static", "The root dir of the static content.")
	metrics   = flag.Bool("metrics", true, "Enable Prometheus metrics.")

	httpRequestsTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	})
)

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)
	if *metrics {
		metricsServer()
	}
	serve(*port, *staticDir)
}

func serve(port string, staticDir string) {
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", handleMiddleware(fs))
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleMiddleware(fs http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpRequestsTotal.Add(1)
		http.StripPrefix("/", fs).ServeHTTP(w, r)
	})
}

func metricsServer() {
	log.Println("registering metrics")
	r := prometheus.NewRegistry()
	r.MustRegister(httpRequestsTotal)
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
}
