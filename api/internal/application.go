package internal

import (
	"github.com/benidevo/vega-ai-landing-page/api/internal/actions"
	"log"
	"net/http"
	"strings"
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
	// Debug logging
	log.Printf("DEBUG: Full URL: %s", r.URL.String())
	log.Printf("DEBUG: Query params: %v", r.URL.Query())
	
	if action := r.URL.Query().Get("action"); action != "" {
		log.Printf("DEBUG: Found action in query: %s", action)
		return action
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	if path != "" {
		log.Printf("DEBUG: Found action in path: %s", path)
		return path
	}

	log.Printf("DEBUG: No action found, returning empty string")
	return ""
}
