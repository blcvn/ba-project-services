package main

import (
	"fmt"
	"log"
	"net"

	"github.com/blcvn/backend/services/feature-service/common/configs"
	"github.com/blcvn/backend/services/feature-service/controllers"
	"github.com/blcvn/backend/services/feature-service/helper"
	repoPostgres "github.com/blcvn/backend/services/feature-service/repository/postgres"
	"github.com/blcvn/backend/services/feature-service/usecases"
	pb "github.com/blcvn/kratos-proto/go/feature"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Load config
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Setup Database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// db.AutoMigrate(&entities.Feature{}) // Should use migration scripts

	// 3. Setup Dependencies
	repo := repoPostgres.NewFeatureRepository(db)
	uc := usecases.NewFeatureUsecase(repo)
	transform := helper.NewTransform()
	controller := controllers.NewFeatureController(uc, transform)

	// 4. Start gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFeatureServiceServer(s, controller)

	log.Printf("Feature Service listening on port %d", cfg.Server.GrpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
