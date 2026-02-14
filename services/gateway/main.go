package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nadeesha_medagama/branch-aware-ci/pkg/interfaces"
)

const (
	defaultPort              = "8080"
	defaultBranchDetectorURL = "http://branch-detector:8081"
	defaultPolicyEngineURL   = "http://policy-engine:8082"
	defaultConfigServiceURL  = "http://config-service:8083"
)

type Gateway struct {
	branchDetectorURL string
	policyEngineURL   string
	configServiceURL  string
	httpClient        *http.Client
}

func main() {
	port := getEnv("PORT", defaultPort)

	gateway := &Gateway{
		branchDetectorURL: getEnv("BRANCH_DETECTOR_URL", defaultBranchDetectorURL),
		policyEngineURL:   getEnv("POLICY_ENGINE_URL", defaultPolicyEngineURL),
		configServiceURL:  getEnv("CONFIG_SERVICE_URL", defaultConfigServiceURL),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/health", gateway.HealthCheck)
	mux.HandleFunc("/ready", gateway.ReadinessCheck)
	mux.HandleFunc("/api/v1/analyze", gateway.AnalyzeBranch)

	// Legacy CLI compatibility endpoint
	mux.HandleFunc("/api/v1/detect-and-decide", gateway.DetectAndDecide)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMiddleware(loggingMiddleware(mux)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("API Gateway listening on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API Gateway...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

// AnalyzeBranch orchestrates the full analysis workflow
func (g *Gateway) AnalyzeBranch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RepoPath   string `json:"repo_path"`
		ConfigPath string `json:"config_path,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Step 1: Detect branch
	branchInfo, err := g.detectBranch(r.Context(), req.RepoPath)
	if err != nil {
		respondWithError(w, fmt.Sprintf("Branch detection failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Step 2: Get configuration
	config, err := g.getConfig(r.Context(), req.ConfigPath)
	if err != nil {
		respondWithError(w, fmt.Sprintf("Config retrieval failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Step 3: Evaluate policy
	decision, err := g.evaluatePolicy(r.Context(), branchInfo, config)
	if err != nil {
		respondWithError(w, fmt.Sprintf("Policy evaluation failed: %v", err), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, decision, http.StatusOK)
}

// DetectAndDecide is a simpler endpoint for CLI compatibility
func (g *Gateway) DetectAndDecide(w http.ResponseWriter, r *http.Request) {
	g.AnalyzeBranch(w, r)
}

// detectBranch calls the branch detector service
func (g *Gateway) detectBranch(ctx context.Context, repoPath string) (*interfaces.BranchInfo, error) {
	reqBody, _ := json.Marshal(map[string]string{"repo_path": repoPath})

	resp, err := g.httpClient.Post(
		g.branchDetectorURL+"/api/v1/detect",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		BranchInfo *interfaces.BranchInfo `json:"branch_info"`
		Error      string                 `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Error != "" {
		return nil, fmt.Errorf(result.Error)
	}

	return result.BranchInfo, nil
}

// getConfig calls the config service
func (g *Gateway) getConfig(ctx context.Context, configPath string) (*interfaces.Config, error) {
	// For now, return default config
	// TODO: Implement actual config service call
	return getDefaultConfig(), nil
}

// evaluatePolicy calls the policy engine service
func (g *Gateway) evaluatePolicy(ctx context.Context, branchInfo *interfaces.BranchInfo, config *interfaces.Config) (*interfaces.Decision, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"branch_info": branchInfo,
		"config":      config,
	})

	resp, err := g.httpClient.Post(
		g.policyEngineURL+"/api/v1/evaluate",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Decision *interfaces.Decision `json:"decision"`
		Error    string               `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Error != "" {
		return nil, fmt.Errorf(result.Error)
	}

	return result.Decision, nil
}

// HealthCheck endpoint
func (g *Gateway) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, map[string]string{"status": "healthy"}, http.StatusOK)
}

// ReadinessCheck endpoint
func (g *Gateway) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	// TODO: Check downstream services
	respondWithJSON(w, map[string]string{"status": "ready"}, http.StatusOK)
}

// Middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDefaultConfig() *interfaces.Config {
	return &interfaces.Config{
		Environments: map[string]interfaces.EnvironmentConfig{
			"production": {
				Name:             "production",
				RequiresApproval: true,
				AllowedBranches:  []string{"main", "master"},
				Variables:        map[string]string{"ENV": "production"},
				NotifyOnDeploy:   true,
			},
			"staging": {
				Name:             "staging",
				RequiresApproval: false,
				AllowedBranches:  []string{"staging", "develop"},
				Variables:        map[string]string{"ENV": "staging"},
				NotifyOnDeploy:   true,
			},
			"development": {
				Name:             "development",
				RequiresApproval: false,
				AllowedBranches:  []string{"feature/*", "bugfix/*"},
				Variables:        map[string]string{"ENV": "development"},
				NotifyOnDeploy:   false,
			},
		},
		BranchMappings: []interfaces.BranchMapping{
			{Pattern: "main", Environment: "production", Actions: []string{"test", "deploy", "notify"}, Priority: 100},
			{Pattern: "master", Environment: "production", Actions: []string{"test", "deploy", "notify"}, Priority: 100},
			{Pattern: "staging", Environment: "staging", Actions: []string{"test", "deploy"}, Priority: 90},
			{Pattern: "feature/*", Environment: "development", Actions: []string{"test"}, Priority: 50},
		},
		Policies: interfaces.PolicyConfig{
			RequireTests:       true,
			RequireCodeReview:  true,
			AutoDeployBranches: []string{"main", "staging"},
		},
	}
}
