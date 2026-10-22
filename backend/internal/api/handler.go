// Package api exposes the calculator over HTTP.
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/GuiSebax/sezzle-calculator/backend/internal/calculator"
)

// maxBodyBytes caps request bodies; a valid request is well under 100 bytes.
const maxBodyBytes = 1 << 10 // 1 KiB

// operation adapts a calculator function to a common signature.
// Unary operations ignore b and do not require it in the request.
type operation struct {
	unary bool
	fn    func(a, b float64) (float64, error)
}

var operations = map[string]operation{
	"add":        {fn: calculator.Add},
	"subtract":   {fn: calculator.Subtract},
	"multiply":   {fn: calculator.Multiply},
	"divide":     {fn: calculator.Divide},
	"power":      {fn: calculator.Power},
	"percentage": {fn: calculator.Percentage},
	"sqrt": {unary: true, fn: func(a, _ float64) (float64, error) {
		return calculator.Sqrt(a)
	}},
}

// calculateRequest uses pointers so a missing field can be told apart from 0.
type calculateRequest struct {
	A *float64 `json:"a"`
	B *float64 `json:"b"`
}

type calculateResponse struct {
	Operation string  `json:"operation"`
	Result    float64 `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// NewRouter returns the service's HTTP handler with CORS applied.
func NewRouter(allowedOrigins []string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /api/v1/calculate/{operation}", handleCalculate)
	return withCORS(mux, allowedOrigins)
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("operation")
	op, ok := operations[name]
	if !ok {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("unknown operation %q", name))
		return
	}

	req, status, err := decodeRequest(w, r)
	if err != nil {
		writeError(w, status, err.Error())
		return
	}
	if req.A == nil {
		writeError(w, http.StatusBadRequest, `field "a" is required`)
		return
	}
	var b float64
	if !op.unary {
		if req.B == nil {
			writeError(w, http.StatusBadRequest, `field "b" is required`)
			return
		}
		b = *req.B
	}

	result, err := op.fn(*req.A, b)
	if err != nil {
		writeError(w, statusForCalcError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, calculateResponse{Operation: name, Result: result})
}

// decodeRequest parses a single JSON object from a size-limited body and
// returns the HTTP status to use if it fails.
func decodeRequest(w http.ResponseWriter, r *http.Request) (calculateRequest, int, error) {
	var req calculateRequest
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))

	if err := dec.Decode(&req); err != nil {
		var maxErr *http.MaxBytesError
		var typeErr *json.UnmarshalTypeError
		switch {
		case errors.As(err, &maxErr):
			return req, http.StatusRequestEntityTooLarge,
				fmt.Errorf("request body must not exceed %d bytes", maxErr.Limit)
		case errors.Is(err, io.EOF):
			return req, http.StatusBadRequest, errors.New("request body must not be empty")
		case errors.As(err, &typeErr) && typeErr.Field == "":
			return req, http.StatusBadRequest, errors.New("request body must be a JSON object")
		case errors.As(err, &typeErr):
			return req, http.StatusBadRequest, fmt.Errorf("field %q must be a number", typeErr.Field)
		default:
			return req, http.StatusBadRequest, errors.New("request body must be valid JSON")
		}
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return req, http.StatusBadRequest, errors.New("request body must contain a single JSON object")
	}
	return req, http.StatusOK, nil
}

// statusForCalcError maps calculator errors to HTTP statuses. Domain errors
// are well-formed requests the math cannot satisfy, hence 422.
func statusForCalcError(err error) int {
	switch {
	case errors.Is(err, calculator.ErrDivisionByZero),
		errors.Is(err, calculator.ErrNegativeSqrt),
		errors.Is(err, calculator.ErrNonFinite):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encode response", "err", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}
