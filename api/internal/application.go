package internal

import (
	"log"
	"net/http"
	"strings"
	"vega.ai/landing-api/internal/actions"
)

// Application is an HTTP handler that manages CORS headers and routes requests
// based on an extracted action from the URL path or query parameters.
func Application(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	action := extractAction(r)
	log.Printf("INFO: Processing request with action: %s from %s", action, r.RemoteAddr)

	switch action {
	case ActionFeedback:
		actions.HandleFeedback(w, r)
	default:
		log.Printf("ERROR: Unknown action requested: %s", action)
		http.Error(w, "Unknown action", http.StatusBadRequest)
	}
}

func extractAction(r *http.Request) string {
	if action := r.URL.Query().Get("action"); action != "" {
		return action
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	if path != "" {
		return path
	}

	return ""
}
