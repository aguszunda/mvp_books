package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"datadog-exercise/internal/controller"
	"datadog-exercise/internal/domain"
	"datadog-exercise/internal/repository"
	"datadog-exercise/internal/service"
	"datadog-exercise/platform"
)

func main() {
	db, err := platform.InitDB()
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

	bookHandler.RegisterRoutes(r)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
