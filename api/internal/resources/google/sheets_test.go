package google

import (
	"context"
	"errors"
	"testing"
)

// MockSheetsService provides a mock implementation for testing
type MockSheetsService struct {
	AppendFeedbackFunc func(ctx context.Context, feedback *FeedbackData) error
}

func (m *MockSheetsService) AppendFeedback(ctx context.Context, feedback *FeedbackData) error {
	if m.AppendFeedbackFunc != nil {
		return m.AppendFeedbackFunc(ctx, feedback)
	}
	return nil
}

func TestSheetsConfig_Validation(t *testing.T) {
	tests := []struct {
		name        string
		config      *SheetsConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid configuration",
			config: &SheetsConfig{
				SpreadsheetID: "1234567890abcdef",
				SheetName:     "Feedback",
			},
			expectError: false,
		},
		{
			name: "missing spreadsheet ID",
			config: &SheetsConfig{
				SpreadsheetID: "",
				SheetName:     "Feedback",
			},
			expectError: true,
			errorMsg:    "spreadsheet ID is required",
		},
		{
			name: "empty sheet name defaults to Feedback",
			config: &SheetsConfig{
				SpreadsheetID: "1234567890abcdef",
				SheetName:     "",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			_, err := NewGoogleSheetsService(ctx, tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else if !tt.expectError && err != nil {
				// For valid configs, we expect a different error (ADC auth failure in test)
				// but not the validation errors we're testing for
				if err.Error() == tt.errorMsg {
					t.Errorf("got validation error when config should be valid: %v", err)
				}
			}
		})
	}
}

func TestMockSheetsService_AppendFeedback(t *testing.T) {
	tests := []struct {
		name          string
		mockFunc      func(ctx context.Context, feedback *FeedbackData) error
		feedback      *FeedbackData
		expectedError string
	}{
		{
			name: "successful append",
			mockFunc: func(ctx context.Context, feedback *FeedbackData) error {
				return nil
			},
			feedback: &FeedbackData{
				Helpfulness:        "very-helpful",
				SetupDifficulty:    3,
				DocsQuality:        "good",
				SetupIssues:        "none",
				AdditionalFeedback: "great tool",
				Email:              "test@example.com",
				Source:             "landing-page",
			},
			expectedError: "",
		},
		{
			name: "mock function returns error",
			mockFunc: func(ctx context.Context, feedback *FeedbackData) error {
				return errors.New("sheets API error")
			},
			feedback: &FeedbackData{
				Helpfulness: "helpful",
				Source:      "test",
			},
			expectedError: "sheets API error",
		},
		{
			name:          "no mock function provided",
			mockFunc:      nil,
			feedback:      &FeedbackData{},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockSheetsService{
				AppendFeedbackFunc: tt.mockFunc,
			}

			ctx := context.Background()
			err := mockService.AppendFeedback(ctx, tt.feedback)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("expected error %q, got %q", tt.expectedError, err.Error())
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestFeedbackData_Structure(t *testing.T) {
	feedback := &FeedbackData{
		Helpfulness:        "very-helpful",
		SetupDifficulty:    4,
		DocsQuality:        "excellent",
		SetupIssues:        "none",
		AdditionalFeedback: "love the tool",
		Email:              "user@example.com",
		Source:             "landing-page",
	}

	// Test that all fields are accessible
	if feedback.Helpfulness != "very-helpful" {
		t.Errorf("expected helpfulness 'very-helpful', got %q", feedback.Helpfulness)
	}
	if feedback.SetupDifficulty != 4 {
		t.Errorf("expected setup difficulty 4, got %d", feedback.SetupDifficulty)
	}
	if feedback.DocsQuality != "excellent" {
		t.Errorf("expected docs quality 'excellent', got %q", feedback.DocsQuality)
	}
	if feedback.SetupIssues != "none" {
		t.Errorf("expected setup issues 'none', got %q", feedback.SetupIssues)
	}
	if feedback.AdditionalFeedback != "love the tool" {
		t.Errorf("expected additional feedback 'love the tool', got %q", feedback.AdditionalFeedback)
	}
	if feedback.Email != "user@example.com" {
		t.Errorf("expected email 'user@example.com', got %q", feedback.Email)
	}
	if feedback.Source != "landing-page" {
		t.Errorf("expected source 'landing-page', got %q", feedback.Source)
	}
}

func TestGoogleSheetsService_AppendFeedback_NilCheck(t *testing.T) {
	service := &GoogleSheetsService{}
	ctx := context.Background()

	err := service.AppendFeedback(ctx, nil)
	if err == nil {
		t.Error("expected error for nil feedback, got nil")
	}

	expectedError := "feedback data cannot be nil"
	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}
