package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"datadog-exercise/internal/adapter/handler"
	"datadog-exercise/internal/adapter/repository"
	"datadog-exercise/internal/domain"
	"datadog-exercise/internal/infrastructure"
	"datadog-exercise/internal/service"
)

func main() {
	// 1. Initialize OpenTelemetry
	shutdown := infrastructure.InitProvider()
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown TracerProvider: %v", err)
		}
	}()

	// 2. Initialize Database
	db, err := infrastructure.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// 3. Run Migrations
	// Note: In production, use a proper migration tool (e.g., golang-migrate)
	if err := db.AutoMigrate(&domain.Book{}); err != nil {
		log.Fatal(err)
	}

	bookRepo := repository.NewMysqlBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)

	r := gin.Default()

	bookHandler.RegisterRoutes(r)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
