package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"mvp_books/internal/controller"
	"mvp_books/internal/domain"
	"mvp_books/internal/middleware"
	"mvp_books/internal/repository"
	"mvp_books/internal/service"
	"mvp_books/platform"
	"mvp_books/platform/connection"
	"mvp_books/platform/logger"
	"mvp_books/platform/metrics"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func initMetrics(meter metric.Meter) error {
	return metrics.Init(meter)
}

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logger.Init(logLevel)

	ctx := context.Background()

	shutdown, err := platform.InitTelemetry(ctx)
	if err != nil {
		slog.Error("failed to initialize telemetry", "error", err)
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			slog.Error("failed to shutdown telemetry", "error", err)
		}
	}()

	meter := otel.GetMeterProvider().Meter("books-api")
	if err := initMetrics(meter); err != nil {
		slog.Error("failed to initialize metrics", "error", err)
		log.Fatal(err)
	}

	db, err := connection.InitDB(ctx)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&domain.Book{}); err != nil {
		slog.Error("failed to migrate database", "error", err)
		log.Fatal(err)
	}

	bookRepo := repository.NewMysqlBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := controller.NewBookHandler(bookService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logging())
	r.Use(middleware.Metrics("books-api"))

	bookHandler.RegisterRoutes(r)

	slog.Info("server starting", "port", "8080")
	if err := r.Run(":8080"); err != nil {
		slog.Error("server failed", "error", err)
		log.Fatal(err)
	}
}
