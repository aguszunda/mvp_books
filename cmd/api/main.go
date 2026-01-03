package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"datadog-exercise/internal/adapter/handler"
	"datadog-exercise/internal/adapter/repository"
	"datadog-exercise/internal/domain"
	"datadog-exercise/internal/infrastructure"
	"datadog-exercise/internal/service"
)

func main() {
	db, err := infrastructure.InitDB()
	if err != nil {
		log.Fatal(err)
	}

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
