package actions

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/benidevo/vega-ai-landing-page/api/internal/resources/google"
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

var (
	sheetsService google.SheetsService
	sheetsOnce    sync.Once
)

// initSheetsService initializes the Google Sheets service once
func initSheetsService() {
	sheetsOnce.Do(func() {
		ctx := context.Background()
		service, err := google.NewGoogleSheetsServiceFromEnv(ctx)
		if err != nil {
			log.Printf("WARNING: Google Sheets not configured: %v", err)
			return
		}

		sheetsService = service
		log.Printf("INFO: Google Sheets service initialized successfully")
	})
}

func HandleFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("ERROR: Invalid method %s for feedback endpoint", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	initSheetsService()

	var req FeedbackRequest
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
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

		setupIssues := strings.Join(r.Form["setupIssues"], ", ")

		req = FeedbackRequest{
			Helpfulness:        r.FormValue("helpfulness"),
			SetupDifficulty:    setupDifficulty,
			DocsQuality:        r.FormValue("docsQuality"),
			SetupIssues:        setupIssues,
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

	log.Printf("INFO: Processing feedback from source:  %s", req.Source)

	if sheetsService != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		feedbackData := &google.FeedbackData{
			Helpfulness:        req.Helpfulness,
			SetupDifficulty:    req.SetupDifficulty,
			DocsQuality:        req.DocsQuality,
			SetupIssues:        req.SetupIssues,
			AdditionalFeedback: req.AdditionalFeedback,
			Email:              req.Email,
			Source:             req.Source,
		}
		if err := sheetsService.AppendFeedback(ctx, feedbackData); err != nil {
			log.Printf("ERROR: Failed to store feedback in Google Sheets: %v", err)
			// continue processing
		}
	} else {
		log.Printf("WARNING: Google Sheets service not available, feedback not stored in sheets")
	}

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
