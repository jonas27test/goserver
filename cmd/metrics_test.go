package main

import (
	"testing"
)

func TestCustomMetrics(t *testing.T) {
	metricsServer()
	unencryptedServer(":8080", "../static")
}
