package handler

import (
	"encoding/json"
	"net/http"

	"github.com/NadeeshaMedagama/branch_aware_ci/pkg/interfaces"
)

// HTTPHandler handles HTTP requests for the policy engine service
type HTTPHandler struct {
	engine interfaces.IPolicyEngine
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(engine interfaces.IPolicyEngine) *HTTPHandler {
	return &HTTPHandler{
		engine: engine,
	}
}

// EvaluatePolicyRequest represents the HTTP request body
type EvaluatePolicyRequest struct {
	BranchInfo *interfaces.BranchInfo `json:"branch_info"`
	Config     *interfaces.Config     `json:"config"`
}

// EvaluatePolicyResponse represents the HTTP response body
type EvaluatePolicyResponse struct {
	Decision *interfaces.Decision `json:"decision"`
	Error    string               `json:"error,omitempty"`
}

// EvaluatePolicy handles policy evaluation requests
func (h *HTTPHandler) EvaluatePolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req EvaluatePolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	decision, err := h.engine.Evaluate(r.Context(), req.BranchInfo, req.Config)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, EvaluatePolicyResponse{
		Decision: decision,
	}, http.StatusOK)
}

// HealthCheck handles health check requests
func (h *HTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, map[string]string{"status": "healthy"}, http.StatusOK)
}

// ReadinessCheck handles readiness check requests
func (h *HTTPHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, map[string]string{"status": "ready"}, http.StatusOK)
}

// Helper functions
func respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	respondWithJSON(w, map[string]string{"error": message}, statusCode)
}
