package actions

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleFeedback_Success(t *testing.T) {
	req := FeedbackRequest{
		Helpfulness:        "excellent",
		SetupDifficulty:    3,
		DocsQuality:        "good",
		SetupIssues:        "none",
		AdditionalFeedback: "Great tool!",
		Email:              "test@example.com",
		Source:             "test",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	HandleFeedback(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response FeedbackResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	if response.Message == "" {
		t.Error("Expected non-empty message")
	}
}

func TestHandleFeedback_InvalidMethod(t *testing.T) {
	httpReq := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	HandleFeedback(w, httpReq)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandleFeedback_InvalidJSON(t *testing.T) {
	httpReq := httptest.NewRequest("POST", "/", strings.NewReader("invalid json"))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	HandleFeedback(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleFeedback_MissingRequiredField(t *testing.T) {
	req := FeedbackRequest{
		// Missing Helpfulness
		SetupDifficulty: 3,
		DocsQuality:     "good",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	HandleFeedback(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleFeedback_DefaultSource(t *testing.T) {
	req := FeedbackRequest{
		Helpfulness: "good",
		// Source is empty, should default to "landing-page"
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	HandleFeedback(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestFeedbackRequest_JSONTags(t *testing.T) {
	req := FeedbackRequest{
		Helpfulness:        "excellent",
		SetupDifficulty:    5,
		DocsQuality:        "good",
		SetupIssues:        "none",
		AdditionalFeedback: "Great!",
		Email:              "test@example.com",
		Source:             "test",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded FeedbackRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.Helpfulness != req.Helpfulness {
		t.Error("Helpfulness not preserved")
	}
	if decoded.SetupDifficulty != req.SetupDifficulty {
		t.Error("SetupDifficulty not preserved")
	}
}
