package controller

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mvp_books/internal/domain"
	"mvp_books/internal/service"
)

type BookHandler struct {
	service service.BookService
}

func NewBookHandler(service service.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var book domain.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		slog.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(c.Request.Context(), &book); err != nil {
		slog.Error("failed to create book", "error", err, "title", book.Title, "author", book.Author)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("book created", "id", book.ID, "title", book.Title, "author", book.Author)
	c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	books, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		slog.Error("failed to list books", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("books listed", "count", len(books))
	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	id := c.Param("id")
	if _, err := strconv.Atoi(id); err != nil {
		slog.Warn("invalid book id", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	book, err := h.service.GetOne(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			slog.Warn("book not found", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		}
		slog.Error("failed to get book", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("book retrieved", "id", book.ID, "title", book.Title)
	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/books", h.CreateBook)
	r.GET("/books", h.GetBooks)
	r.GET("/books/:id", h.GetBook)
}
