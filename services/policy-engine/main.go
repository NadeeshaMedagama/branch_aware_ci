package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NadeeshaMedagama/branch_aware_ci/pkg/interfaces"
	"github.com/NadeeshaMedagama/branch_aware_ci/services/policy-engine/engine"
	"github.com/NadeeshaMedagama/branch_aware_ci/services/policy-engine/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	defaultGRPCPort = "50052"
	defaultHTTPPort = "8082"
)

func main() {
	grpcPort := getEnv("GRPC_PORT", defaultGRPCPort)
	httpPort := getEnv("HTTP_PORT", defaultHTTPPort)

	// Create policy engine instance
	var policyEngine interfaces.IPolicyEngine = engine.NewPolicyEngine()

	// Start gRPC server
	go startGRPCServer(grpcPort, policyEngine)

	// Start HTTP server
	go startHTTPServer(httpPort, policyEngine)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Policy Engine Service...")
}

func startGRPCServer(port string, engine interfaces.IPolicyEngine) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	log.Printf("Policy Engine gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func startHTTPServer(port string, engine interfaces.IPolicyEngine) {
	h := handler.NewHTTPHandler(engine)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/ready", h.ReadinessCheck)
	mux.HandleFunc("/api/v1/evaluate", h.EvaluatePolicy)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Policy Engine HTTP server listening on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
