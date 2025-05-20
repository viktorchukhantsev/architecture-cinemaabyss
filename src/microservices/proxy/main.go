package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	"proxy/internal/strangler"
)

const (
	envPort                   = "PORT"
	envMonolithURL            = "MONOLITH_URL"
	envMoviesURL              = "MOVIES_SERVICE_URL"
	envEventsURL              = "EVENTS_SERVICE_URL"
	envGradualMigration       = "GRADUAL_MIGRATION"
	envMoviesMigrationPercent = "MOVIES_MIGRATION_PERCENT"

	defaultPort = "8000"
)

func main() {
	monolithURL, err := url.Parse(os.Getenv(envMonolithURL))
	if err != nil {
		log.Fatal("Unable to parse monolith URL")
	}

	eventsURL, err := url.Parse(os.Getenv(envEventsURL))
	if err != nil {
		log.Fatal("Unable to parse events URL")
	}

	_, err = url.Parse(os.Getenv(envMoviesURL))
	if err != nil {
		log.Fatal("Unable to parse movies URL")
	}

	percentageStr := os.Getenv(envMoviesMigrationPercent)
	percentage, err := strconv.Atoi(percentageStr)
	if err != nil {
		log.Fatal("Unable to parse movies migration percentage")
	}

	proxyMonolith := httputil.NewSingleHostReverseProxy(monolithURL)
	eventsProxy := httputil.NewSingleHostReverseProxy(eventsURL)

	http.HandleFunc("/health", handleHealth)
	http.Handle("/api/events/health", eventsProxy)

	if os.Getenv(envGradualMigration) == "true" {
		moviesStranglerProxy := strangler.NewStranglerRouter(
			os.Getenv(envMonolithURL),
			os.Getenv(envMoviesURL),
			percentage,
		)

		http.HandleFunc("/api/movies", moviesStranglerProxy.ServeHTTP)
	} else {
		http.Handle("/api/movies", proxyMonolith)
	}

	http.Handle("/api/users", proxyMonolith)
	http.Handle("/api/payments", proxyMonolith)
	http.Handle("/api/subscription", proxyMonolith)
	http.Handle("/api/events/movie", eventsProxy)
	http.Handle("/api/events/user", eventsProxy)
	http.Handle("/api/events/payment", eventsProxy)

	port := os.Getenv(envPort)
	if port == "" {
		port = defaultPort
	}
	log.Printf("Starting proxy microservice on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}
