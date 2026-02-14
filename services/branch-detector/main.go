package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nadeesha_medagama/branch-aware-ci/pkg/interfaces"
	"github.com/nadeesha_medagama/branch-aware-ci/services/branch-detector/detector"
	"github.com/nadeesha_medagama/branch-aware-ci/services/branch-detector/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	defaultGRPCPort = "50051"
	defaultHTTPPort = "8081"
)

func main() {
	// Get ports from environment or use defaults
	grpcPort := getEnv("GRPC_PORT", defaultGRPCPort)
	httpPort := getEnv("HTTP_PORT", defaultHTTPPort)

	// Create detector instance (following Dependency Inversion principle)
	var branchDetector interfaces.IBranchDetector = detector.NewBranchDetector()

	// Start gRPC server
	go startGRPCServer(grpcPort, branchDetector)

	// Start HTTP server for health checks
	go startHTTPServer(httpPort, branchDetector)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Branch Detector Service...")
}

func startGRPCServer(port string, detector interfaces.IBranchDetector) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	// Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection service (for development/debugging)
	reflection.Register(grpcServer)

	// TODO: Register BranchDetectorService when proto is generated

	log.Printf("Branch Detector gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func startHTTPServer(port string, detector interfaces.IBranchDetector) {
	h := handler.NewHTTPHandler(detector)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/ready", h.ReadinessCheck)
	mux.HandleFunc("/api/v1/detect", h.DetectBranch)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Branch Detector HTTP server listening on port %s", port)
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
