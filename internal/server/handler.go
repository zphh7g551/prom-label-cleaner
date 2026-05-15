package server

import (
	"log"
	"net/http"

	"github.com/prom-label-cleaner/internal/exporter"
	"github.com/prom-label-cleaner/internal/pipeline"
)

// MetricsHandler returns an http.HandlerFunc that runs the pipeline
// and writes the cleaned metrics to the response.
func MetricsHandler(runner *pipeline.Runner, exp *exporter.Exporter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := runner.Run(r.Context())
		if err != nil {
			log.Printf("pipeline run error: %v", err)
			http.Error(w, "failed to scrape metrics", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

		if err := exp.Write(result.Families); err != nil {
			log.Printf("exporter write error: %v", err)
			http.Error(w, "failed to write metrics", http.StatusInternalServerError)
			return
		}
	}
}

// HealthHandler returns a simple liveness probe handler.
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}
