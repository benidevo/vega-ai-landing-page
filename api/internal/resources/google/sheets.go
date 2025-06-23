package google

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// SheetsService defines the interface for Google Sheets operations
type SheetsService interface {
	AppendFeedback(ctx context.Context, feedback *FeedbackData) error
}

// FeedbackData represents feedback data for storage in Google Sheets
type FeedbackData struct {
	Helpfulness        string
	SetupDifficulty    int
	DocsQuality        string
	SetupIssues        string
	AdditionalFeedback string
	Email              string
	Source             string
}

// GoogleSheetsService handles Google Sheets operations
type GoogleSheetsService struct {
	service       *sheets.Service
	spreadsheetID string
	sheetName     string
}

// SheetsConfig holds configuration for Google Sheets service
type SheetsConfig struct {
	SpreadsheetID string
	SheetName     string
}

// NewGoogleSheetsServiceFromEnv creates a new Google Sheets service from environment variables
func NewGoogleSheetsServiceFromEnv(ctx context.Context) (SheetsService, error) {
	spreadsheetID := os.Getenv("GOOGLE_SPREADSHEET_ID")
	sheetName := os.Getenv("GOOGLE_SHEET_NAME")

	if spreadsheetID == "" {
		return nil, fmt.Errorf("GOOGLE_SPREADSHEET_ID environment variable is required")
	}

	config := &SheetsConfig{
		SpreadsheetID: spreadsheetID,
		SheetName:     sheetName,
	}

	return NewGoogleSheetsService(ctx, config)
}

// NewGoogleSheetsService creates a new Google Sheets service instance
func NewGoogleSheetsService(ctx context.Context, config *SheetsConfig) (*GoogleSheetsService, error) {
	if config.SpreadsheetID == "" {
		return nil, fmt.Errorf("spreadsheet ID is required")
	}
	if config.SheetName == "" {
		config.SheetName = "Vega AI Feedback" // default sheet name
	}

	// Try multiple authentication methods
	var service *sheets.Service
	var err error

	// Method 1: Try with impersonation if service account is specified
	serviceAccountEmail := os.Getenv("GCP_SERVICE_ACCOUNT_EMAIL")
	if serviceAccountEmail == "" {
		serviceAccountEmail = "vega-feedback-sheets@vega-ai-live.iam.gserviceaccount.com"
	}

	// First, try to get default credentials
	creds, err := google.FindDefaultCredentials(ctx, sheets.SpreadsheetsScope)
	if err != nil {
		log.Printf("WARNING: Could not find default credentials: %v", err)
		// Try without credentials (will use metadata service in Cloud Run)
		service, err = sheets.NewService(ctx, option.WithScopes(sheets.SpreadsheetsScope))
	} else {
		// If we have credentials, try impersonation
		ts, impersonateErr := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
			TargetPrincipal: serviceAccountEmail,
			Scopes:          []string{sheets.SpreadsheetsScope},
			Delegates:       []string{},
		}, option.WithCredentials(creds))

		if impersonateErr != nil {
			log.Printf("WARNING: Impersonation failed, falling back to default credentials: %v", impersonateErr)
			// Fall back to using the credentials directly
			service, err = sheets.NewService(ctx,
				option.WithCredentials(creds),
				option.WithScopes(sheets.SpreadsheetsScope),
			)
			if err != nil {
				log.Printf("ERROR: Failed to create sheets service with default credentials - Type: %T, Error: %v", err, err)
				return nil, fmt.Errorf("failed to create sheets service with default credentials: %w", err)
			}
			log.Printf("INFO: Successfully created sheets service with default credentials")
		} else {
			// Use impersonated credentials
			service, err = sheets.NewService(ctx, option.WithTokenSource(ts))
			if err != nil {
				log.Printf("ERROR: Failed to create sheets service with impersonation - Type: %T, Error: %v", err, err)
				return nil, fmt.Errorf("failed to create sheets service with impersonation: %w", err)
			}
			log.Printf("INFO: Successfully created sheets service with impersonation for account")
		}
	}

	if err != nil {
		log.Printf("ERROR: Failed to create sheets service - Type: %T, Error: %v", err, err)
		return nil, fmt.Errorf("failed to create sheets service: %w", err)
	}

	sheetsService := &GoogleSheetsService{
		service:       service,
		spreadsheetID: config.SpreadsheetID,
		sheetName:     config.SheetName,
	}

	if err := sheetsService.ensureHeaders(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize sheet headers: %w", err)
	}

	log.Printf("INFO: Google Sheets service initialized successfully for spreadsheet: %s", config.SpreadsheetID)
	return sheetsService, nil
}

// AppendFeedback appends feedback data to the Google Sheet
func (g *GoogleSheetsService) AppendFeedback(ctx context.Context, feedback *FeedbackData) error {
	if feedback == nil {
		return fmt.Errorf("feedback data cannot be nil")
	}

	values := []any{
		time.Now().Format(time.RFC3339),
		feedback.Helpfulness,
		feedback.SetupDifficulty,
		feedback.DocsQuality,
		feedback.SetupIssues,
		feedback.AdditionalFeedback,
		feedback.Email,
		feedback.Source,
	}

	valueRange := &sheets.ValueRange{
		Values: [][]any{values},
	}

	range_ := fmt.Sprintf("'%s'!A:H", g.sheetName)
	appendCall := g.service.Spreadsheets.Values.Append(g.spreadsheetID, range_, valueRange)
	appendCall.ValueInputOption("RAW")
	appendCall.InsertDataOption("INSERT_ROWS")

	_, err := appendCall.Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to append feedback to sheet: %w", err)
	}

	log.Printf("INFO: Successfully appended feedback to Google Sheets from source: %s", feedback.Source)
	return nil
}

// ensureHeaders ensures the sheet has proper headers
func (g *GoogleSheetsService) ensureHeaders(ctx context.Context) error {
	// Sheet names with spaces must be wrapped in single quotes
	range_ := fmt.Sprintf("'%s'!A1:H1", g.sheetName)
	response, err := g.service.Spreadsheets.Values.Get(g.spreadsheetID, range_).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to check existing headers: %w", err)
	}

	// If no data or empty first row, add headers
	if len(response.Values) == 0 || len(response.Values[0]) == 0 {
		headers := []interface{}{
			"Timestamp",
			"Helpfulness",
			"Setup Difficulty",
			"Docs Quality",
			"Setup Issues",
			"Additional Feedback",
			"Email",
			"Source",
		}

		valueRange := &sheets.ValueRange{
			Values: [][]any{headers},
		}

		updateCall := g.service.Spreadsheets.Values.Update(g.spreadsheetID, range_, valueRange)
		updateCall.ValueInputOption("RAW")

		_, err := updateCall.Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("failed to add headers: %w", err)
		}

		log.Printf("INFO: Added headers to Google Sheet: %s", g.sheetName)
	}

	return nil
}
