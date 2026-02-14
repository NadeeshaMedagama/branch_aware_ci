package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nadeesha_medagama/branch-aware-ci/pkg/interfaces"
)

// HTTPHandler handles HTTP requests for the branch detector service
// Following Interface Segregation Principle: depends only on what it needs
type HTTPHandler struct {
	detector interfaces.IBranchDetector
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(detector interfaces.IBranchDetector) *HTTPHandler {
	return &HTTPHandler{
		detector: detector,
	}
}

// DetectBranchRequest represents the HTTP request body
type DetectBranchRequest struct {
	RepoPath string `json:"repo_path"`
}

// DetectBranchResponse represents the HTTP response body
type DetectBranchResponse struct {
	BranchInfo *interfaces.BranchInfo `json:"branch_info"`
	Error      string                 `json:"error,omitempty"`
}

// DetectBranch handles branch detection requests
func (h *HTTPHandler) DetectBranch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DetectBranchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	branchInfo, err := h.detector.DetectBranch(r.Context(), req.RepoPath)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, DetectBranchResponse{
		BranchInfo: branchInfo,
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
