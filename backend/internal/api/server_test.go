package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCORSPreflight(t *testing.T) {
	// Set up environment for testing
	oldOrigin := os.Getenv("ALLOWED_ORIGIN")
	os.Setenv("ALLOWED_ORIGIN", "http://34.58.122.79")
	defer os.Setenv("ALLOWED_ORIGIN", oldOrigin)

	s := New()

	tests := []struct {
		name           string
		origin         string
		method         string
		wantStatus     int
		wantAllowOrigin string
		wantVary       string
	}{
		{
			name:           "OPTIONS request from allowed origin",
			origin:         "http://34.58.122.79",
			method:         http.MethodOptions,
			wantStatus:     http.StatusNoContent,
			wantAllowOrigin: "http://34.58.122.79",
			wantVary:       "Origin",
		},
		{
			name:           "OPTIONS request from different origin",
			origin:         "http://different.com",
			method:         http.MethodOptions,
			wantStatus:     http.StatusNoContent,
			wantAllowOrigin: "",
			wantVary:       "",
		},
		{
			name:           "POST request from allowed origin with error response",
			origin:         "http://34.58.122.79",
			method:         http.MethodPost,
			wantStatus:     http.StatusBadRequest, // expected due to missing body
			wantAllowOrigin: "http://34.58.122.79",
			wantVary:       "Origin",
		},
		{
			name:           "GET request from allowed origin with success response",
			origin:         "http://34.58.122.79",
			method:         http.MethodGet,
			wantStatus:     http.StatusOK,
			wantAllowOrigin: "http://34.58.122.79",
			wantVary:       "Origin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use /health endpoint for GET tests, /api/v1/evaluate for POST/OPTIONS tests
			path := "/api/v1/evaluate"
			if tt.method == http.MethodGet {
				path = "/health"
			}
			req := httptest.NewRequest(tt.method, path, nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			rec := httptest.NewRecorder()

			s.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if got := rec.Header().Get("Access-Control-Allow-Origin"); got != tt.wantAllowOrigin {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, tt.wantAllowOrigin)
			}

			if got := rec.Header().Get("Vary"); got != tt.wantVary {
				t.Errorf("Vary = %q, want %q", got, tt.wantVary)
			}

			// For allowed origins, check other CORS headers are present
			if tt.wantAllowOrigin != "" {
				if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
					t.Error("Access-Control-Allow-Methods should be set for allowed origin")
				}
				if got := rec.Header().Get("Access-Control-Allow-Headers"); got == "" {
					t.Error("Access-Control-Allow-Headers should be set for allowed origin")
				}
			}

			// Verify CORS headers are set even on error responses
			// This test specifically checks that non-OPTIONS requests from allowed origins
			// have CORS headers set, which browsers need on all responses including errors
			if tt.method == http.MethodPost && tt.wantAllowOrigin != "" && rec.Code >= 400 {
				if rec.Header().Get("Access-Control-Allow-Origin") == "" {
					t.Error("CORS headers must be present on error responses for allowed origins")
				}
			}
		})
	}
}
