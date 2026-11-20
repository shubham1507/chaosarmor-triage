package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const apiAddress = "127.0.0.1:8080"

func runAPIServer() error {
	server := &http.Server{
		Addr:              apiAddress,
		Handler:           newAPIRouter(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Starting ChaosArmor Triage API on http://%s", server.Addr)
	return server.ListenAndServe()
}

func newAPIRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/healthz", handleHealth)
	return router
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprintln(w, `{"status":"ok"}`); err != nil {
		log.Printf("write health response: %v", err)
	}
}
