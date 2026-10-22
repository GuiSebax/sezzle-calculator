package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCORS(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		origin      string
		preflight   bool
		wantStatus  int
		wantAllowed string
	}{
		{name: "allowed origin", method: http.MethodPost, origin: "http://localhost:5173", wantStatus: http.StatusOK, wantAllowed: "http://localhost:5173"},
		{name: "disallowed origin", method: http.MethodPost, origin: "http://evil.example", wantStatus: http.StatusOK, wantAllowed: ""},
		{name: "no origin", method: http.MethodPost, wantStatus: http.StatusOK, wantAllowed: ""},
		{name: "preflight allowed", method: http.MethodOptions, origin: "http://localhost:5173", preflight: true, wantStatus: http.StatusNoContent, wantAllowed: "http://localhost:5173"},
		{name: "preflight disallowed", method: http.MethodOptions, origin: "http://evil.example", preflight: true, wantStatus: http.StatusNoContent, wantAllowed: ""},
	}

	router := NewRouter(testOrigins)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/v1/calculate/add", strings.NewReader(`{"a": 1, "b": 2}`))
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.preflight {
				req.Header.Set("Access-Control-Request-Method", http.MethodPost)
			}
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := rec.Header().Get("Access-Control-Allow-Origin"); got != tt.wantAllowed {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, tt.wantAllowed)
			}
			if got := rec.Header().Get("Vary"); got != "Origin" {
				t.Errorf("Vary = %q, want Origin", got)
			}
			if tt.preflight && rec.Header().Get("Access-Control-Allow-Methods") == "" {
				t.Error("preflight response missing Access-Control-Allow-Methods")
			}
		})
	}
}
