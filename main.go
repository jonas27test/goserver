package main

import (
	"flag"
	"log"
	"net/http"
	"path"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	port       = flag.String("port", ":8080", "The ort")
	serverCert = flag.String("serverCrt", "tls.crt", "The location of the certificate. Must be below cert dir.")
	serverKey  = flag.String("serverKey", "tls.key", "The location of the key. Must be below cert dir")
	metrics    = flag.Bool("metrics", true, "Enable Prometheus metrics.")
	secure     = flag.Bool("secure", false, "Should the server use TLS.")

	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"code", "method"})

	httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Duration of all HTTP requests",
	}, []string{"code", "handler", "method"})
)

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)
	certDir := "./cert"
	*serverCert, *serverKey = path.Join(certDir, *serverCert), path.Join(certDir, *serverKey)
	if *metrics {
		metricsServer()
	}
	if *secure {
		log.Println("Start TLS server")
		tlsServer(*serverCert, *serverKey)
	} else {
		log.Println("Start unencrypted server")
		unencryptedServer()
	}
}

func tlsServer(crt string, key string) {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", http.StripPrefix("/", fs))
	err := http.ListenAndServeTLS(*port, *serverCert, *serverKey, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func unencryptedServer() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", http.StripPrefix("/", fs))
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func metricsServer() {
	log.Println("registering metrics")
	// r := prometheus.NewRegistry()
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	http.Handle("/metrics", promhttp.Handler())
}
