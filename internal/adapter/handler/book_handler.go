package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"datadog-exercise/internal/domain"
	"datadog-exercise/internal/port"
)

type BookHandler struct {
	service port.BookService
}

func NewBookHandler(service port.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var book domain.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the current span (Tracing logic preserved from original)
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("book.title", book.Title))

	if err := h.service.Create(c.Request.Context(), &book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	// Simulate latency as in original code
	time.Sleep(50 * time.Millisecond)

	books, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	id := c.Param("id")
	book, err := h.service.GetOne(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/books", h.CreateBook)
	r.GET("/books", h.GetBooks)
	r.GET("/books/:id", h.GetBook)
}
