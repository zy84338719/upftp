package response

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccess(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Test Success function
	data := map[string]interface{}{
		"message": "test",
	}
	Success(w, data)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestError(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Test Error function
	Error(w, http.StatusInternalServerError, "test error")

	// Check response status code
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestBadRequest(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Test BadRequest function
	BadRequest(w, "test error")

	// Check response status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUnauthorized(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Test Unauthorized function
	Unauthorized(w, "test error")

	// Check response status code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestForbidden(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Test Forbidden function
	Forbidden(w, "test error")

	// Check response status code
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestNotFound(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Test NotFound function
	NotFound(w, "test error")

	// Check response status code
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}
