package main

import (
	"log"
	"net/http"
	"os"

	promptsRepo "github.com/blcvn/backend/services/prompt-service/repository/postgres"
	"github.com/blcvn/backend/services/prompt-service/usecases"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize Database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost user=postgres password=postgres dbname=ba_agent port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Initialize Repository
	repo := promptsRepo.NewPromptRepository(db)

	// Initialize Usecase
	uc := usecases.NewPromptUsecase(repo)
	_ = uc

	// Setup Routes (stub)
	http.HandleFunc("/prompts", func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here using uc.ListPrompts() etc.
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	log.Println("Prompt Service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
