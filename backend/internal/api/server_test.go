package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSPreflight(t *testing.T) {
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
			name:           "POST request from allowed origin",
			origin:         "http://34.58.122.79",
			method:         http.MethodPost,
			wantStatus:     http.StatusBadRequest, // will fail due to no body, but that's ok
			wantAllowOrigin: "http://34.58.122.79",
			wantVary:       "Origin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/v1/evaluate", nil)
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

			// For allowed origins, check other CORS headers
			if tt.wantAllowOrigin != "" && tt.method == http.MethodOptions {
				if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
					t.Error("Access-Control-Allow-Methods should be set")
				}
				if got := rec.Header().Get("Access-Control-Allow-Headers"); got == "" {
					t.Error("Access-Control-Allow-Headers should be set")
				}
			}
		})
	}
}
