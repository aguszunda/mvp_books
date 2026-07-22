package main

import (
	"context"
	"log"

	"mvp_books/internal/controller"
	"mvp_books/internal/domain"
	"mvp_books/internal/middleware"
	"mvp_books/internal/repository"
	"mvp_books/internal/service"
	"mvp_books/platform"
	"mvp_books/platform/connection"
	"mvp_books/platform/metrics"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func initMetrics(meter metric.Meter) error {
	return metrics.Init(meter)
}

func main() {
	ctx := context.Background()

	shutdown, err := platform.InitTelemetry(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown telemetry: %v", err)
		}
	}()

	meter := otel.GetMeterProvider().Meter("books-api")
	if err := initMetrics(meter); err != nil {
		log.Fatalf("Failed to initialize metrics: %v", err)
	}

	db, err := connection.InitDB(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&domain.Book{}); err != nil {
		log.Fatal(err)
	}

	bookRepo := repository.NewMysqlBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := controller.NewBookHandler(bookService)

	r := gin.Default()
	r.Use(middleware.Metrics("books-api"))

	bookHandler.RegisterRoutes(r)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
