package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApplication_CORS(t *testing.T) {
	req := httptest.NewRequest("OPTIONS", "/", nil)
	w := httptest.NewRecorder()

	Application(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS, got %d", w.Code)
	}

	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type",
	}

	for header, expected := range expectedHeaders {
		if got := w.Header().Get(header); got != expected {
			t.Errorf("Expected %s header to be %s, got %s", header, expected, got)
		}
	}
}

func TestApplication_UnknownAction(t *testing.T) {
	req := httptest.NewRequest("POST", "/?action=unknown", nil)
	w := httptest.NewRecorder()

	Application(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for unknown action, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Unknown action") {
		t.Error("Expected 'Unknown action' in response body")
	}
}

func TestApplication_FeedbackAction(t *testing.T) {
	// Test with query parameter
	req := httptest.NewRequest("POST", "/?action=feedback", strings.NewReader(`{"helpfulness":"good"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Application(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for feedback action, got %d", w.Code)
	}
}

func TestApplication_FeedbackActionFromPath(t *testing.T) {
	// Test with path
	req := httptest.NewRequest("POST", "/feedback", strings.NewReader(`{"helpfulness":"good"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Application(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for feedback action from path, got %d", w.Code)
	}
}

func TestExtractAction_QueryParameter(t *testing.T) {
	req := httptest.NewRequest("GET", "/?action=feedback", nil)
	action := extractAction(req)

	if action != "feedback" {
		t.Errorf("Expected action 'feedback', got '%s'", action)
	}
}

func TestExtractAction_Path(t *testing.T) {
	req := httptest.NewRequest("GET", "/feedback", nil)
	action := extractAction(req)

	if action != "feedback" {
		t.Errorf("Expected action 'feedback', got '%s'", action)
	}
}

func TestExtractAction_Empty(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	action := extractAction(req)

	if action != "" {
		t.Errorf("Expected empty action, got '%s'", action)
	}
}

func TestExtractAction_QueryParameterPriority(t *testing.T) {
	// Query parameter should take priority over path
	req := httptest.NewRequest("GET", "/somepath?action=feedback", nil)
	action := extractAction(req)

	if action != "feedback" {
		t.Errorf("Expected action 'feedback' from query parameter, got '%s'", action)
	}
}
