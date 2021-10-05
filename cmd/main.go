package main

import (
	"compress/gzip"
	"flag"
	"io"
	"log"
	"net/http"
	"strings"

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
	log.Println("starting server on port: " + *port)
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

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// Use the Writer part of gzipResponseWriter to write the output.

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func handleMiddleware(fs http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			httpRequestsTotal.Add(1)
			http.StripPrefix("/", fs).ServeHTTP(w, r)
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		http.StripPrefix("/", fs).ServeHTTP(gzipResponseWriter{gz, w}, r)
	})
}

func metricsServer() {
	log.Println("registering metrics")
	r := prometheus.NewRegistry()
	r.MustRegister(httpRequestsTotal)
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
}
