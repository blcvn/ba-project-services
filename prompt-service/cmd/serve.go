package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/blcvn/backend/services/prompt-service/controllers"
	"github.com/blcvn/backend/services/prompt-service/repository/postgres"
	"github.com/blcvn/backend/services/prompt-service/usecases"
	pb "github.com/blcvn/kratos-proto/go/prompt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pgDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Prompt Service",
	Run:   runServe,
}

func runServe(cmd *cobra.Command, args []string) {
	dbURL := getEnv("DATABASE_URL", "postgresql://baagent:password@localhost:5432/baagent?sslmode=disable")
	grpcPort := getEnv("GRPC_PORT", "9086")
	httpPort := getEnv("HTTP_PORT", "8086")

	db, err := gorm.Open(pgDriver.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	repo := postgres.NewPromptRepository(db)
	usecase := usecases.NewPromptUsecase(repo)
	controller := controllers.NewPromptController(usecase)

	grpcServer := grpc.NewServer()
	pb.RegisterPromptServiceServer(grpcServer, controller)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed on grpc port: %v", err)
	}

	go func() {
		log.Printf("Starting gRPC on %s", grpcPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("gRPC serve error: %v", err)
		}
	}()

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = pb.RegisterPromptServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%s", grpcPort), opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", httpPort),
		Handler: mux,
	}

	go func() {
		log.Printf("Starting HTTP on %s", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP serve error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	grpcServer.GracefulStop()
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
