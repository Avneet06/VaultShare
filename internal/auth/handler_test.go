package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	reqBody := `{"email":"simple@example.com","password":"123456"}`
	req := httptest.NewRequest("TEST", "/register", strings.NewReader(reqBody)) 

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	RegisterHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}
}
