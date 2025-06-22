package actions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// FeedbackRequest represents the structure of feedback submitted by users.
type FeedbackRequest struct {
	Helpfulness        string `json:"helpfulness"`
	SetupDifficulty    int    `json:"setupDifficulty"`
	DocsQuality        string `json:"docsQuality"`
	SetupIssues        string `json:"setupIssues"`
	AdditionalFeedback string `json:"additionalFeedback"`
	Email              string `json:"email"`
	Source             string `json:"source"`
}

// FeedbackResponse represents the standard response structure for feedback-related API endpoints.
type FeedbackResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func HandleFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("ERROR: Invalid method %s for feedback endpoint", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FeedbackRequest
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		// Handle JSON request (from manual fetch calls)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("ERROR: Failed to decode JSON request: %v", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			log.Printf("ERROR: Failed to parse form data: %v", err)
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		// Convert form data to struct
		setupDifficulty := 5 // default value
		if val := r.FormValue("setupDifficulty"); val != "" {
			if parsed, err := strconv.Atoi(val); err == nil {
				setupDifficulty = parsed
			}
		}

		req = FeedbackRequest{
			Helpfulness:        r.FormValue("helpfulness"),
			SetupDifficulty:    setupDifficulty,
			DocsQuality:        r.FormValue("docsQuality"),
			SetupIssues:        r.FormValue("setupIssues"),
			AdditionalFeedback: r.FormValue("additionalFeedback"),
			Email:              r.FormValue("email"),
			Source:             r.FormValue("source"),
		}
	}

	if req.Helpfulness == "" {
		log.Printf("ERROR: Missing required field 'helpfulness' in feedback request")
		http.Error(w, "Helpfulness is required", http.StatusBadRequest)
		return
	}

	if req.Source == "" {
		req.Source = "landing-page"
	}

	log.Printf("INFO: Processing feedback from source: %s, email: %s", req.Source, req.Email)

	// TODO: Store feedback in Google Sheets
	fmt.Printf("Feedback received at %s: %+v\n", time.Now().Format(time.RFC3339), req)

	log.Printf("INFO: Feedback processed successfully")

	w.Header().Set("Content-Type", "application/json")
	response := FeedbackResponse{
		Success: true,
		Message: "Thank you for your feedback! Your insights will help us improve Vega AI for everyone.",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("ERROR: Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
