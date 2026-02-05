package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/blcvn/backend/services/prompt-service/controllers"
	"github.com/blcvn/backend/services/prompt-service/helper"
	"github.com/blcvn/backend/services/prompt-service/usecases"
	"github.com/blcvn/kratos-proto/go/prompt"

	postgres_repo "github.com/blcvn/backend/services/prompt-service/repository/postgres"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Config
	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	// Database
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Init Layers
	repo := postgres_repo.NewPromptRepository(db)
	tmplEngine := helper.NewTemplateEngine()
	usecase := usecases.NewPromptUsecase(repo, tmplEngine)
	controller := controllers.NewPromptController(usecase)

	// gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	prompt.RegisterPromptServiceServer(s, controller)

	log.Printf("Prompt Service listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
