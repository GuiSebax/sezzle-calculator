package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testOrigins = []string{"http://localhost:5173"}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name       string
		operation  string
		body       string
		wantStatus int
		wantResult float64 // checked when wantStatus is 200
		wantError  string  // checked when wantStatus is not 200
	}{
		{name: "add", operation: "add", body: `{"a": 5, "b": 3}`, wantStatus: http.StatusOK, wantResult: 8},
		{name: "subtract", operation: "subtract", body: `{"a": 5, "b": 3}`, wantStatus: http.StatusOK, wantResult: 2},
		{name: "multiply", operation: "multiply", body: `{"a": 5, "b": 3}`, wantStatus: http.StatusOK, wantResult: 15},
		{name: "divide", operation: "divide", body: `{"a": 6, "b": 3}`, wantStatus: http.StatusOK, wantResult: 2},
		{name: "power", operation: "power", body: `{"a": 2, "b": 10}`, wantStatus: http.StatusOK, wantResult: 1024},
		{name: "sqrt with only a", operation: "sqrt", body: `{"a": 16}`, wantStatus: http.StatusOK, wantResult: 4},
		{name: "sqrt ignores b", operation: "sqrt", body: `{"a": 9, "b": 100}`, wantStatus: http.StatusOK, wantResult: 3},
		{name: "percentage", operation: "percentage", body: `{"a": 10, "b": 50}`, wantStatus: http.StatusOK, wantResult: 5},
		{name: "zero is a valid value", operation: "add", body: `{"a": 0, "b": 0}`, wantStatus: http.StatusOK, wantResult: 0},

		{name: "unknown operation", operation: "modulo", body: `{"a": 5, "b": 3}`, wantStatus: http.StatusBadRequest, wantError: `unknown operation "modulo"`},
		{name: "missing a", operation: "add", body: `{"b": 3}`, wantStatus: http.StatusBadRequest, wantError: `field "a" is required`},
		{name: "missing b", operation: "add", body: `{"a": 5}`, wantStatus: http.StatusBadRequest, wantError: `field "b" is required`},
		{name: "null b", operation: "add", body: `{"a": 5, "b": null}`, wantStatus: http.StatusBadRequest, wantError: `field "b" is required`},
		{name: "sqrt missing a", operation: "sqrt", body: `{}`, wantStatus: http.StatusBadRequest, wantError: `field "a" is required`},
		{name: "empty body", operation: "add", body: ``, wantStatus: http.StatusBadRequest, wantError: "request body must not be empty"},
		{name: "malformed JSON", operation: "add", body: `{"a": 5,`, wantStatus: http.StatusBadRequest, wantError: "request body must be valid JSON"},
		{name: "wrong type", operation: "add", body: `{"a": "5", "b": 3}`, wantStatus: http.StatusBadRequest, wantError: `field "a" must be a number`},
		{name: "not an object", operation: "add", body: `[1, 2]`, wantStatus: http.StatusBadRequest, wantError: "request body must be a JSON object"},
		{name: "trailing data", operation: "add", body: `{"a": 5, "b": 3}{}`, wantStatus: http.StatusBadRequest, wantError: "request body must contain a single JSON object"},
		{name: "body too large", operation: "add", body: `{"a": 5, "b": 3, "pad": "` + strings.Repeat("x", maxBodyBytes) + `"}`, wantStatus: http.StatusRequestEntityTooLarge, wantError: "request body must not exceed 1024 bytes"},

		{name: "division by zero", operation: "divide", body: `{"a": 5, "b": 0}`, wantStatus: http.StatusUnprocessableEntity, wantError: "division by zero"},
		{name: "sqrt of negative", operation: "sqrt", body: `{"a": -4}`, wantStatus: http.StatusUnprocessableEntity, wantError: "square root of a negative number"},
		{name: "non-finite result", operation: "power", body: `{"a": 10, "b": 400}`, wantStatus: http.StatusUnprocessableEntity, wantError: "result is not a finite number"},
	}

	router := NewRouter(testOrigins)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate/"+tt.operation, strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}

			if tt.wantStatus == http.StatusOK {
				var got calculateResponse
				decode(t, rec, &got)
				if got.Operation != tt.operation || got.Result != tt.wantResult {
					t.Errorf("response = %+v, want {Operation:%s Result:%v}", got, tt.operation, tt.wantResult)
				}
				return
			}
			var got errorResponse
			decode(t, rec, &got)
			if got.Error != tt.wantError {
				t.Errorf("error = %q, want %q", got.Error, tt.wantError)
			}
		})
	}
}

func TestCalculateRejectsOtherMethods(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate/add", nil)
	rec := httptest.NewRecorder()

	NewRouter(testOrigins).ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	NewRouter(testOrigins).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var got map[string]string
	decode(t, rec, &got)
	if got["status"] != "ok" {
		t.Errorf("status field = %q, want ok", got["status"])
	}
}

func decode(t *testing.T, rec *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.NewDecoder(rec.Body).Decode(v); err != nil {
		t.Fatalf("decode response %q: %v", rec.Body, err)
	}
}
