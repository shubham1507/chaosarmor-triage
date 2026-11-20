package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIRouter(t *testing.T) {
	router := newAPIRouter()
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{"health check", http.MethodGet, "/healthz", http.StatusOK},
		{"unsupported method", http.MethodPost, "/healthz", http.StatusMethodNotAllowed},
		{"unknown path", http.MethodGet, "/missing", http.StatusNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}

			switch test.wantStatus {
			case http.StatusOK:
				if got := response.Header().Get("Content-Type"); got != "application/json" {
					t.Errorf("Content-Type = %q, want application/json", got)
				}
				if got := response.Body.String(); got != "{\"status\":\"ok\"}\n" {
					t.Errorf("body = %q, want health JSON", got)
				}
			case http.StatusMethodNotAllowed:
				if got := response.Header().Get("Allow"); got != http.MethodGet {
					t.Errorf("Allow = %q, want GET", got)
				}
			}
		})
	}
}
